package auth

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/cedar-policy/cedar-go"
	"github.com/cedar-policy/cedar-go/types"
)

var policyData []byte

type Authorizer struct { //evaluator
	policySet *cedar.PolicySet
}

func NewAuthorizer() (*Authorizer, error) {
	policySet, err := cedar.NewPolicyListFromBytes("policies.cedar", policyData)
	if err != nil {
		return nil, fmt.Errorf("failed to complie cedar policies %w", err)
	}

	log.Println("Successfully complied Cedar RBAC policies!")
	return &Authorizer{policySet: policySet}, nil

}
func (a *Authorizer) IsAuthorized(userID string, role string, actionName string, resourceType string, resourceID string) (bool, error) {
	principalUID := types.NewEntityUID("DevTrack::User", types.String(userID))
	principalRecord := types.Record{
		"role": types.String(role),
	}
	principalEntity := types.Entity{
		UID:        principalUID,
		Attributes: principalRecord,
	}
	actionUID := types.NewEntityUID("DevTrack::Action", types.String(actionName))

	resourceUID:=types.NewEntityUID(types.EntityType("DevTrack::"+resourceType),types.String(resourceID))
	resourceEntity:=types.Entity{
		UID:resourceUID,
	}
	entities:=types.EntitiesFromSlice([]types.Entity{
		principalEntity,
		resourceEntity
	})

	req:=cedar.Request{
		Principal: principalUID,
		Action: actionUID,
		Resource: resourceUID,
	}

	decision,_:=a.policySet.IsAuthorized(entities,req)
	return bool(decision)==bool(cedar.Allow),nil
}
