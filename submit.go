package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type SubmissionBody struct {
	ContactEmail string `json:"contact_email"`
	GithubUrl    string `json:"github_url"`
}

func generateTOTP(secret string, unixTime int64, digits int) string {
	h := hmac.New(sha512.New, []byte(secret))
	b := make([]byte, 8)

	counter := uint64(math.Floor(float64(unixTime) / float64(30)))
	binary.BigEndian.PutUint64(b, counter)

	h.Write(b)
	hash := h.Sum(nil)

	// get the first 4 bits of last byte of hash
	offset := hash[len(hash)-1] & 0xf
	var res int64
	res = int64((int(hash[offset])&0x7f)<<24) |
		int64((int(hash[offset+1])&0xff)<<16) |
		int64((int(hash[offset+2])&0xff)<<8) |
		int64(int(hash[offset+3]) & 0xff)

	otp := int32(res) % int32(math.Pow10(digits))
	format := "%0" + fmt.Sprintf("%d", digits) + "d"
	return fmt.Sprintf(format, otp)
}

func main() {

	data := SubmissionBody{
		ContactEmail: "rsmohamad@ust.hk",
		GithubUrl:    "https://gist.github.com/rsmohamad/17d8af84190d26be03bb6a5d37b68316"}

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&data)

	secret := fmt.Sprintf("%sHDECHALLENGE003", data.ContactEmail)
	password := generateTOTP(secret, time.Now().Unix(), 10)

	req, _ := http.NewRequest(http.MethodPost, "https://hdechallenge-solve.appspot.com/challenge/003/endpoint", body)
	req.SetBasicAuth(data.ContactEmail, password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	//requestDump, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(requestDump))

	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}
