package model

import (
    "fmt"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "bytes"
    "io/ioutil"
    "os"
)

type Profile struct {
    Name string             `json:"name"`
    Birthday string         `json:"birthday"`
    Relationship string     `json:"relationship"`
    Note string             `json:"note"`
    LastViewed string       `json:"last_viewed"`
}

type getLovedOnesResponse struct {
    LovedOnes []int `json:"loved_ones"`
}

type LovedOne struct {
    Prof   Profile `json:"profile"`
    Images []string  `json:"images"`
}

type insertDBResponse struct {
    Id int `json:"id"`
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

func (l *LovedOne) InsertIntoDB(userId string) (int, error) {
    b, err := json.Marshal(l.Prof)
    if err != nil {
        return -1, err
    }

    buf := bytes.NewBuffer(b)
    req, err := http.NewRequest("POST","http://localhost/api/v1/users/loved-one", buf)
    if err != nil {
        return -1, err
    }

    req.Header.Set("Authorization", userId)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return -1, err
    }
    defer resp.Body.Close()

    var r insertDBResponse
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return -1, err
    }
    return r.Id, nil
}

func GetIdsOfLovedOnes(userId int) ([]int, error) {
    resp, err := http.Get(fmt.Sprintf("http://localhost/api/v1/users/loved-one?user_id=%d", userId))
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
