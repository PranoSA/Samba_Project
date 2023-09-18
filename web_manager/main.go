package main

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

func main() {

	//S3_Endpoint := os.Getenv("S3_ENDPOINT")
	S3_Client_id := os.Getenv("S3_CLIENT_ID")
	S3_Access_Key := os.Getenv("S3_ACCESS_KEY")

	// Initialize minio client object.
	minioClient, err := minio.NewWithOptions("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(S3_Client_id, S3_Access_Key, ""),
		Secure: false,
	})

	if err != nil {

		log.Fatalln(err)
	}

	var randomtempdirbytes []byte = make([]byte, 32)
	rand.Read(randomtempdirbytes)

	randomtempdirstring := base32.StdEncoding.EncodeToString(randomtempdirbytes)

	//err = minioClient.FGetObject("backend", "/1"+"/"+"my.mp4", "./random.mp4", minio.GetObjectOptions{})
	/*err = minioClient.FGetObject("tre", "Plans_To_Coalesce.png", "./random.png", minio.GetObjectOptions{})
	if err != nil {
		log.Fatal("IDK YOU SUCK ASSS")
	}
	*/

	var shareid string = "12512"
	var fileName string = "bob/fred/ted/mug"
	//Read Into Directory Here
	//var videoIn string = ""

	os.Mkdir(randomtempdirstring, 0777)

	exec.Command("sh", "./dash.sh", "./big_buck_bunny.mp4", randomtempdirstring).Output()

	err = filepath.WalkDir(randomtempdirstring, func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path, d.Name())
		name := shareid + "/" + fileName + "/" + d.Name()
		_, err = minioClient.FPutObject("frontend", name, randomtempdirstring+"/"+d.Name(), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println(err)
		}
		err = os.Remove(randomtempdirstring + "/" + d.Name())
		if err != nil {
			fmt.Print("Failed ToRemove \n")
		}
		return nil
	})

	_, err = minioClient.FPutObject("frontend", shareid+"/"+fileName+"/"+"index.html", "index.html", minio.PutObjectOptions{})

	if err != nil {
		fmt.Println("Failed to put to Minio")
	}

	err = os.Remove(randomtempdirstring)
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}

	//log.Printf("%#v\n", minioClient) // minioClient is now setup

}
