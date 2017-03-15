package model

import (
    "encoding/base64"
    "encoding/json"
    "io/ioutil"
    "os/exec"
    "os"
    "net/http"
    "bytes"
)

type profile struct {
    Name string             `json:"name"`
    Birthday string         `json:"birthday"`
    Relationship string     `json:"relationship"`
    Note string             `json:"note"`
    LastViewed string       `json:"last_viewed"`
}

type LovedOne struct {
    Profile profile `json:"profile"`
    Images []string  `json:"images"`
    UserId int       `json:"user_id"`
}

type insertDBResponse struct {
    Id string `json:"id"`
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
    err := json.NewDecoder(resp.Body).Decode(&r)
    if err = nil {
        return -1, err
    }
    return r.Id, nil
}
