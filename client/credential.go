package internal

import (
	"fmt"
	"github.com/danieljoos/wincred"
	"strings"
)

func getAllCredentials() ([]*wincred.Credential, error) {
	creds, err := wincred.List()
	if err != nil {
		return nil, err
	}

	ret := make([]*wincred.Credential, 0)
	for _, cred := range creds {
		if strings.Index(cred.TargetName, "fse:") == 0 {
			ret = append(ret, cred)
		}
	}

	return ret, nil
}

func addCredential(credName string, serverAddr string, token []byte) error {
	cred := wincred.NewGenericCredential("fse:" + credName)
	cred.UserName = serverAddr
	cred.CredentialBlob = token
	return cred.Write()
}

func removeCredential(credName string) error {
	matched, err := lookupCredential(credName)
	if err != nil {
		return err
	}

	if cred, err := wincred.GetGenericCredential(matched[0].TargetName); err != nil {
		return err
	} else {
		return cred.Delete()
	}
}

func lookupCredential(credName string) (matched []*wincred.Credential, err error) {
	s := "fse:" + credName

	creds, err := getAllCredentials()
	if err != nil {
		return
	}

	matched = make([]*wincred.Credential, 0)
	for i := range creds {
		if strings.Index(creds[i].TargetName, s) == 0 {
			matched = append(matched, creds[i])
		}
	}

	if len(matched) != 1 {
		if len(matched) == 0 {
			err = fmt.Errorf("no credential has specified name or prefix")
		} else {
			err = fmt.Errorf("multiple credentials match provided name")
		}
	}

	return
}
