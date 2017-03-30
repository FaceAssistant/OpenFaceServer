package web

import (
    "net/http"
    "fa/model"
    "fa/s3util"
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
        } else {
            http.Error(w, "Error getting user id", http.StatusInternalServerError)
            return
        }
    }
}

