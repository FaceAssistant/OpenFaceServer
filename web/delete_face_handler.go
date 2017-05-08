package web

import (
    "net/http"
    "fa/model"
    "fa/s3util"
    "fa/openface"
    "io/ioutil"
    "fmt"
    "os"
)

func DeleteFaceHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if uid, ok := r.Context().Value("uid").(string); ok {
            err := model.DeleteLovedOne(r.FormValue("id"), r.Header.Get("Authorization"))
            if err != nil {
                http.Error(w, "Error deleting loved one from database.", http.StatusInternalServerError)
                return
            }

            err = s3util.DeleteFeatures(r.FormValue("id"), uid)
            if err != nil {
                http.Error(w, "Failed to delete S3 object. Err: " + err.Error(), http.StatusInternalServerError)
                return
            }

            dir, err := ioutil.TempDir("/tmp/", uid)
            if err != nil {
                http.Error(w, "Failed to make dir", http.StatusInternalServerError)
                return
            }
            defer os.RemoveAll(dir)

            labelsf, err := os.Create(dir + "/labels.csv")
            if err != nil {
                http.Error(w, "Error making labels.csv", http.StatusInternalServerError)
                return
            }

            repsf, err := os.Create(dir + "/reps.csv")
            if err != nil {
                http.Error(w, "Error making reps.csv", http.StatusInternalServerError)
                return
            }

            err = s3util.GetFeature("labels.csv", uid, r.Header.Get("Authorization"), labelsf)
            if err != nil {
                http.Error(w, "Error getting labels from AWS: " + err.Error(), http.StatusInternalServerError)
                return
            }

            err = s3util.GetFeature("reps.csv", uid, r.Header.Get("Authorization"), repsf)
            if err != nil {
                http.Error(w, "Error getting reps from AWS: " + err.Error(), http.StatusInternalServerError)
                return
            }

            err = openface.CreatePickle(dir)
            if err != nil {
                http.Error(w, "Error creating model: " + err.Error(), http.StatusInternalServerError)
                return
            }

            err = s3util.UploadFile(dir + "/classifier.pkl", fmt.Sprintf("features/%s/classifier.pkl", uid))
            if err != nil {
                http.Error(w, "Error uploading classifier: " + err.Error(), http.StatusInternalServerError)
                return
            }
        } else {
            http.Error(w, "Error getting user id", http.StatusInternalServerError)
            return
        }
    }
}

