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

type profile struct {
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
    Profile profile `json:"profile"`
    Images []string  `json:"images"`
    UserId int       `json:"user_id"`
}

type insertDBResponse struct {
    Id int `json:"id"`
}

func(l *LovedOne) WriteImagesToFile(dir string) error {
    for _, i := range(l.Images) {
        data, err := base64.StdEncoding.DecodeString(i)
        if err != nil {
            return err
        }

        _, err = writeBytesToFile(data, dir)
        if err != nil {
            return err
        }
    }
    return nil
}

func (l *LovedOne) InsertIntoDB() (int, error) {
    b, err := json.Marshal(l.Profile)
    if err != nil {
        return -1, err
    }

    buf := bytes.NewBuffer(b)
    resp, err := http.Post("http://127.0.0.1/api/v1/users/loved-ones", "application/json", buf)
    defer resp.Body.Close()
    if err != nil {
        return -1, err
    }

    var r insertDBResponse
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return -1, err
    }
    return r.Id, nil
}

func GetIdsOfLovedOnes(userId int) ([]int, error) {
    resp, err := http.Get(fmt.Sprintf("http://127.0.0.1/api/v1/users/loved-ones?user_id=%d", userId))
    defer resp.Body.Close()
    if err != nil {
        return nil, err
    }

    var r getLovedOnesResponse
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return nil, err
    }

    return r.LovedOnes, nil
}


func writeBytesToFile(b []byte, dir string) (string, error) {
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
