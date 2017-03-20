package s3util

import (
    "io"
    "os"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "fa/model"
)

func getObject(bucketName string, object string) (io.ReadCloser, error) {
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

func GetClassifier(userId string, dst io.Writer) error {
    var objReader io.ReadCloser
    objReader, err := getObject("faceassist", fmt.Sprintf("features/%s/classifier.pkl", userId))
    if err != nil {
        if awsErr, ok := err.(awserr.Error); ok {
            if awsErr.Code() != "NoSuchKey" {
                return err
            } else {
                objReader, err = getObject("faceassist", "celeb/classifier.pkl")
                if err != nil {
                    return err
                }
                defer objReader.Close()
            }
        } else {
            return err
        }
    }
    defer objReader.Close()

    _, err = io.Copy(dst, objReader)
    if err != nil {
        return err
    }

    return nil
}

func GetFeature(fileName string, userId int, dst io.Writer) error {
    lovedOnes, err := model.GetIdsOfLovedOnes(userId)
    if err != nil{
        return err
    }

    celeb, err := getObject("faceassist", fmt.Sprintf("celeb/%s", fileName))
    if err != nil {
        return err
    }
    defer celeb.Close()

    _, err = io.Copy(dst, celeb)
    if err != nil {
        return err
    }

    for _, id := range lovedOnes {
        objectKey := fmt.Sprintf("features/%d/%d/%s", userId, id, fileName)
        objReader, err := getObject("faceassist", objectKey)
        if err != nil {
            return err
        }

        _, err = io.Copy(dst, objReader)
        if err != nil {
            return err
        }
        objReader.Close()
    }
    return nil
}

func putS3Object(bucketName string, objectKey string, file io.ReadSeeker) error {
    sess, err := session.NewSession()
    if err != nil {
        return err
    }

    svc := s3.New(sess)
    params := &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key: aws.String(objectKey),
        Body: file,
    }

    _, err = svc.PutObject(params)
    if err != nil {
        return err
    }
    //log response
    return nil
}

func UploadFile(file string, objectKey string) error {
    f, err := os.Open(file)
    defer f.Close()
    if err != nil {
        return err
    }
    return putS3Object("faceassist", objectKey, f)
}
