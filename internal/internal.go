package internal

import (
	"fmt"
	"log"
)

func CheckError(err error, msg string) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("%s err: [%v]", msg, err)
	log.Println(err.Error())
	return err
}
