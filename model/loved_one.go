package model

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
    "bytes"
    "io/ioutil"
    "os"
)

type Profile struct {
    Id string               `json:"id"`
    Name string             `json:"name"`
    Birthday string         `json:"birthday"`
    Relationship string     `json:"relationship"`
    Note string             `json:"note"`
    LastViewed string       `json:"last_viewed"`
}

type getLovedOnesResponse struct {
    LovedOnes []string `json:"loved_ones"`
}

type LovedOne struct {
    Prof   Profile `json:"profile"`
    Images []string  `json:"images"`
}

type insertDBResponse struct {
    Id string `json:"id"`
}

func(l *LovedOne) WriteImagesToFile(dir string) error {
    err := os.Mkdir(dir, 0777)
    if err != nil {
        return err
    }

    for _, i := range(l.Images) {
        data, err := base64.StdEncoding.DecodeString(i)
        if err != nil {
            return err
        }

        _, err = WriteBytesToFile(data, dir)
        if err != nil {
            return err
        }
    }
    return nil
}

func (l *LovedOne) InsertIntoDB(idToken string) (string, error) {
    b, err := json.Marshal(l.Prof)
    if err != nil {
        return "", err
    }

    buf := bytes.NewBuffer(b)
    req, err := http.NewRequest("POST","http://localhost/api/v1/users/loved-one", buf)
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", idToken)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var r insertDBResponse
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return "", err
    }
    return r.Id, nil
}

func GetIdsOfLovedOnes(idToken string) ([]string, error) {
    req, err := http.NewRequest("GET", "http://localhost/api/v1/users/loved-one", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", idToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    var r getLovedOnesResponse
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return nil, err
    }

    return r.LovedOnes, nil
}

func GetLovedOneById(id string, idToken string) (*http.Response, error) {
    req, err := http.NewRequest("GET", "http://localhost/api/v1/users/loved-one?id="+id, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", idToken)

    client := &http.Client{}
    return client.Do(req)
}

//Probably moved this to another package
func WriteBytesToFile(b []byte, dir string) (string, error) {
    f, err := ioutil.TempFile(dir, "img-")
    if err != nil {
        return "", err
    }

    if _, err := f.Write(b); err != nil {
        return "", err
    }

    os.Rename(f.Name(), f.Name() + ".jpg")

    defer func() {
        if err = f.Close(); err != nil {
            panic(err)
        }
    }()

    return f.Name() + ".jpg", nil
}

func DeleteLovedOne(id string, idToken string) error {
    req, err := http.NewRequest("DELETE", "http://localhost/api/v1/users/loved-one?id="+id, nil)

    if err != nil {
        return err
    }
    req.Header.Set("Authorization", idToken)

    client := &http.Client{}
    _, err = client.Do(req)
    if err != nil {
        return err
    }
    return nil
}
