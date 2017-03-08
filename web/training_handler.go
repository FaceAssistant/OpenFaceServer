package web

import (
    "net/http"
    "fa/model"
    "fa/openface"
    "encoding/json"
    "os"
)



func TrainingHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var l model.LovedOne
        err := json.NewDecoder(r.Body).Decode(&l)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        dir := "/tmp/" + l.UserId
        err := os.Mkdir(dir, os.ModeTemporary)
        defer  os.RemoveAll(dir)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        imgDir := dir + "images"
        alignDir := dir + "align"
        featureDir := dir + "feature"

        err = l.WriteImagesToFile(imgDir)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = openface.Train(imgDir, alignDir, featureDir)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}
