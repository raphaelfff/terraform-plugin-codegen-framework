// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"

	generatorschema "github.com/raphaelfff/terraform-plugin-codegen-framework/internal/schema"
)

func NewSchemas(spec spec.Specification) (map[string]generatorschema.GeneratorSchema, error) {
	providerSchemas := make(map[string]generatorschema.GeneratorSchema, 1)

	providerSchema, err := NewSchema(spec.Provider)

	if err != nil {
		return nil, err
	}

	providerSchemas[spec.Provider.Name] = providerSchema

	return providerSchemas, nil
}

func NewSchema(p *provider.Provider) (generatorschema.GeneratorSchema, error) {
	var s generatorschema.GeneratorSchema

	if p.Schema == nil {
		return s, nil
	}

	attributes := make(generatorschema.GeneratorAttributes, len(p.Schema.Attributes))
	blocks := make(generatorschema.GeneratorBlocks, len(p.Schema.Blocks))

	for _, v := range p.Schema.Attributes {
		a, err := NewAttribute("", v)

		if err != nil {
			return s, err
		}

		attributes[v.Name] = a
	}

	s.Attributes = attributes

	for _, v := range p.Schema.Blocks {
		b, err := NewBlock(v)

		if err != nil {
			return s, err
		}

		blocks[v.Name] = b
	}

	s.Blocks = blocks

	s.Description = p.Schema.Description

	s.MarkdownDescription = p.Schema.MarkdownDescription

	s.DeprecationMessage = p.Schema.DeprecationMessage

	return s, nil
}

func NewAttributes(name string, a provider.Attributes) (generatorschema.GeneratorAttributes, error) {
	attributes := make(generatorschema.GeneratorAttributes, len(a))

	for _, v := range a {
		attribute, err := NewAttribute(name, v)

		if err != nil {
			return generatorschema.GeneratorAttributes{}, err
		}

		attributes[v.Name] = attribute
	}

	return attributes, nil
}

func NewAttribute(name string, a provider.Attribute) (generatorschema.GeneratorAttribute, error) {
	name = name + "_" + a.Name

	switch {
	case a.Bool != nil:
		return NewGeneratorBoolAttribute(name, a.Bool)
	case a.Float64 != nil:
		return NewGeneratorFloat64Attribute(name, a.Float64)
	case a.Int64 != nil:
		return NewGeneratorInt64Attribute(name, a.Int64)
	case a.List != nil:
		return NewGeneratorListAttribute(name, a.List)
	case a.ListNested != nil:
		return NewGeneratorListNestedAttribute(name, a.ListNested)
	case a.Map != nil:
		return NewGeneratorMapAttribute(name, a.Map)
	case a.MapNested != nil:
		return NewGeneratorMapNestedAttribute(name, a.MapNested)
	case a.Number != nil:
		return NewGeneratorNumberAttribute(name, a.Number)
	case a.Object != nil:
		return NewGeneratorObjectAttribute(name, a.Object)
	case a.Set != nil:
		return NewGeneratorSetAttribute(name, a.Set)
	case a.SetNested != nil:
		return NewGeneratorSetNestedAttribute(name, a.SetNested)
	case a.SingleNested != nil:
		return NewGeneratorSingleNestedAttribute(name, a.SingleNested)
	case a.String != nil:
		return NewGeneratorStringAttribute(name, a.String)
	}

	return nil, fmt.Errorf("attribute type not defined: %+v", a)
}

func NewBlocks(b provider.Blocks) (generatorschema.GeneratorBlocks, error) {
	blocks := make(generatorschema.GeneratorBlocks, len(b))

	for _, v := range b {
		block, err := NewBlock(v)

		if err != nil {
			return generatorschema.GeneratorBlocks{}, err
		}

		blocks[v.Name] = block
	}

	return blocks, nil
}

func NewBlock(b provider.Block) (generatorschema.GeneratorBlock, error) {
	switch {
	case b.ListNested != nil:
		return NewGeneratorListNestedBlock(b.Name, b.ListNested)
	case b.SetNested != nil:
		return NewGeneratorSetNestedBlock(b.Name, b.SetNested)
	case b.SingleNested != nil:
		return NewGeneratorSingleNestedBlock(b.Name, b.SingleNested)
	}

	return nil, fmt.Errorf("block type not defined: %+v", b)
}
