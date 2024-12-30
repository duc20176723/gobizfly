package gobizfly

import "C"
import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	simpleStorageKeyPath = "/key"
)


var _ SimpleStorageKey = (*cloudSimpleStorageKeyService)(nil)

type SimpleStorageKey interface {
	Create(ctx context.Context, s3cr *KeyCreateRequest) (*KeyHaveSercret, error)
	Get(ctx context.Context, id string) (*KeyHaveSercret, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, opts *ListOptions) ([]*KeyInList, error)
}

type cloudSimpleStorageKeyService struct {
	client *Client
}

type KeyCreateRequest struct {
	SubuserId string `json:"subuser_id"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type KeyHaveSercret struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type KeyInList struct {
	User      string `json:"user"`
	AccessKey string `json:"access_key"`
}


func (c *cloudSimpleStorageService) SimpleStorageKey() *cloudSimpleStorageKeyService {
	return &cloudSimpleStorageKeyService{client: c.client}
}

func (c cloudSimpleStorageKeyService) Create(ctx context.Context, dataCreatekey *KeyCreateRequest) (*KeyHaveSercret, error) {
	req, err := c.client.NewRequest(ctx, http.MethodPost, simpleStorageServiceName, c.resourcePath(), &dataCreatekey)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData struct {
		Message string          `json:"message"`
		Key     *KeyHaveSercret `json:"Key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	return respData.Key, nil
}


func (c *cloudSimpleStorageKeyService) resourcePath() string {
	return simpleStorageKeyPath
}

func (c *cloudSimpleStorageKeyService) itemPath(id string) string {
	if id == "" {
		return simpleStorageKeyPath
	}
	return strings.Join([]string{simpleStorageKeyPath, id}, "/")
}


func (c *cloudSimpleStorageKeyService) Delete(ctx context.Context, id string) error {
	req, err := c.client.NewRequest(ctx, http.MethodDelete, simpleStorageServiceName, c.itemPath(id), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return resp.Body.Close()
}


func (c *cloudSimpleStorageKeyService) Get(ctx context.Context, id string) (*KeyHaveSercret, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, simpleStorageServiceName, c.itemPath(id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	key := &KeyHaveSercret{}
	if err := json.NewDecoder(resp.Body).Decode(key); err != nil {
		return nil, err
	}
	return key, nil
}


func (c *cloudSimpleStorageKeyService) List(ctx context.Context, opts *ListOptions) ([]*KeyInList, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, simpleStorageServiceName, c.resourcePath(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Keys []*KeyInList `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Keys, nil
}
