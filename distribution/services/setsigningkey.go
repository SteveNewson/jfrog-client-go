package services

import (
	"encoding/json"
	"net/http"

	artifactoryUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SetSigningKeyService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
}

func NewSetSigningKeyService(client *jfroghttpclient.JfrogHttpClient) *SetSigningKeyService {
	return &SetSigningKeyService{client: client}
}

func (ssk *SetSigningKeyService) GetDistDetails() auth.ServiceDetails {
	return ssk.DistDetails
}

func (ssk *SetSigningKeyService) SetSigningKey(signBundleParams SetSigningKeyParams) error {
	body := &SetSigningKeyBody{
		PublicKey:  signBundleParams.PublicKey,
		PrivateKey: signBundleParams.PrivateKey,
	}
	return ssk.execSetSigningKey(body)
}

func (ssk *SetSigningKeyService) execSetSigningKey(setSigningKeyBody *SetSigningKeyBody) error {
	httpClientsDetails := ssk.DistDetails.CreateHttpClientDetails()
	content, err := json.Marshal(setSigningKeyBody)
	if err != nil {
		return errorutils.CheckError(err)
	}
	url := ssk.DistDetails.GetUrl() + "/api/v1/keys/pgp"
	artifactoryUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ssk.client.SendPut(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	return errorutils.CheckError(err)
}

type SetSigningKeyBody struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type SetSigningKeyParams struct {
	PublicKey  string
	PrivateKey string
}

func NewSetSigningKeyParams(publicKey, privateKey string) SetSigningKeyParams {
	return SetSigningKeyParams{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}
