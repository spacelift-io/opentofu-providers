package tf5

import (
	"context"
	"fmt"

	"github.com/apparentlymart/opentofu-providers/tofuprovider/grpc/tfplugin5"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/internal/common"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerops"
	"github.com/apparentlymart/opentofu-providers/tofuprovider/providerschema"
)

// ApplyManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ApplyManagedResourceChange(ctx context.Context, req *providerops.ApplyManagedResourceChangeRequest) (providerops.ApplyManagedResourceChangeResponse, error) {
	panic("unimplemented")
}

// ImportManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ImportManagedResourceState(ctx context.Context, req *providerops.ImportManagedResourceStateRequest) (providerops.ImportManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// MoveManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) MoveManagedResourceState(ctx context.Context, req *providerops.MoveManagedResourceStateRequest) (providerops.MoveManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// PlanManagedResourceChange implements tofuprovider.GRPCPluginProvider.
func (p *Provider) PlanManagedResourceChange(ctx context.Context, req *providerops.PlanManagedResourceChangeRequest) (providerops.PlanManagedResourceChangeResponse, error) {
	priorState, err := makeDynamicValueMsgpack(req.PriorState)
	if err != nil {
		return nil, fmt.Errorf("invalid PriorState value: %w", err)
	}
	proposedNewState, err := makeDynamicValueMsgpack(req.ProposedNewState)
	if err != nil {
		return nil, fmt.Errorf("invalid ProposedNewState value: %w", err)
	}
	config, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}

	var providerMeta *tfplugin5.DynamicValue
	if req.ProviderMeta != providerschema.NoDynamicValue {
		providerMeta, err = makeDynamicValueMsgpack(req.ProviderMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid ProviderMeta value: %w", err)
		}
	}

	protoReq := &tfplugin5.PlanResourceChange_Request{
		TypeName:           req.ResourceType,
		PriorState:         priorState,
		ProposedNewState:   proposedNewState,
		Config:             config,
		PriorPrivate:       req.PriorProviderInternal,
		ProviderMeta:       providerMeta,
		ClientCapabilities: prepareClientCapabilities(req.ClientCapabilities),
	}

	protoResp, err := p.client.PlanResourceChange(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return planManagedResourceChangeResponse{proto: protoResp}, nil
}

// ReadManagedResource implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ReadManagedResource(ctx context.Context, req *providerops.ReadManagedResourceRequest) (providerops.ReadManagedResourceResponse, error) {
	panic("unimplemented")
}

// UpgradeManagedResourceState implements tofuprovider.GRPCPluginProvider.
func (p *Provider) UpgradeManagedResourceState(ctx context.Context, req *providerops.UpgradeManagedResourceStateRequest) (providerops.UpgradeManagedResourceStateResponse, error) {
	panic("unimplemented")
}

// ValidateManagedResourceConfig implements tofuprovider.GRPCPluginProvider.
func (p *Provider) ValidateManagedResourceConfig(ctx context.Context, req *providerops.ValidateManagedResourceConfigRequest) (providerops.ValidateManagedResourceConfigResponse, error) {
	configVal, err := makeDynamicValueMsgpack(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid Config value: %w", err)
	}
	protoReq := &tfplugin5.ValidateResourceTypeConfig_Request{
		TypeName: req.ResourceType,
		Config:   configVal,
	}

	protoResp, err := p.client.ValidateResourceTypeConfig(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return validateManagedResourceConfigResponse{proto: protoResp}, nil
}

type validateManagedResourceConfigResponse struct {
	proto *tfplugin5.ValidateResourceTypeConfig_Response
	common.SealedImpl
}

// Diagnostics implements providerops.ValidateEphemeralResourceConfigResponse.
func (v validateManagedResourceConfigResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: v.proto.Diagnostics}
}

type planManagedResourceChangeResponse struct {
	proto *tfplugin5.PlanResourceChange_Response
	common.SealedImpl
}

// Diagnostics implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) Diagnostics() providerops.Diagnostics {
	return diagnostics{proto: p.proto.Diagnostics}
}

// PlannedNewState implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) PlannedNewState() providerschema.DynamicValueOut {
	return dynamicValue{proto: p.proto.PlannedState}
}

// PlannedProviderInternal implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) PlannedProviderInternal() []byte {
	return p.proto.PlannedPrivate
}

// LegacyTypeSystem implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) LegacyTypeSystem() bool {
	return p.proto.LegacyTypeSystem
}

// Deferred implements providerops.PlanManagedResourceChangeResponse.
func (p planManagedResourceChangeResponse) Deferred() providerops.Deferred {
	if p.proto.Deferred == nil {
		return nil
	}
	return deferred{proto: p.proto.Deferred}
}
