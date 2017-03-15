package s3utils

import (
    "io"
    "ioutil"
    "bufio"
    "os"
    "fmt"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "net/http"
    "encoding/json"
)

type getLovedOnesResponse struct {
    LovedOnes []int `json:"loved_ones"`
}

func GetObject(bucketName string, object string) (io.ReadCloser, error) {
    var r io.ReadCloser
    sess, err := session.NewSession()
    if err != nil {
        return r, err
    }

    svc := s3.New(sess)

    params := &s3.GetObjectInput {
        Bucket: aws.String(bucketName),
        Key: aws.String(object),
    }

    resp, err := svc.GetObject(params)
    if err != nil {
        return r, err
    }

    return resp.Body, nil
}

func GetLabels(userId string) (*File, error) {
    //Refactor
    file, err := ioutil.TempFile("/tmp/", userId)
    if err != nil {
        return "", err
    }

    fileWriter := bufio.NewWriter(file)

    resp, err := http.Get("http://127.0.0.1/api/v1/users/loved-ones?user_id=" + userId)
    if err != nil {
        return "", err
    }

    var r getLovedOnesResponse
    err := json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
        return "", err
    }

    for _, id := range r.LovedOnes {
        objReader, err := GetObject("fa", fmt.Sprintf("features/%d/%d/labels.csv", userId, id))
        if err != nil {
            return "", err
        }

        //FINISH THIS
        _, err := io.Copy(writer, objReader)
        if err != nil {
            return "", err
        }
    }

    err := os.Rename(file.Name(), file.Name() + ".csv")
    if err != nil {
        return "", err
    }
    //Get list of profile ids from fa-db.
    //Loop through each id, create urls for s3 get
    //make get request, concat bytes to one file
    //return file
}
