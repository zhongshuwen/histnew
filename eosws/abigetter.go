// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eosws

import (
	"context"
	"fmt"

	pbstatedb "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/statedb/v1"
	"github.com/zhongshuwen/zswchain-go"
)

type ABIGetter interface {
	GetABI(ctx context.Context, blockNum uint32, account zsw.AccountName) (*zsw.ABI, error)
}

type DefaultABIGetter struct {
	client pbstatedb.StateClient
}

func NewDefaultABIGetter(client pbstatedb.StateClient) *DefaultABIGetter {
	return &DefaultABIGetter{
		client: client,
	}
}

func (g *DefaultABIGetter) GetABI(ctx context.Context, blockNum uint32, contract zsw.AccountName) (*zsw.ABI, error) {
	response, err := g.client.GetABI(ctx, &pbstatedb.GetABIRequest{BlockNum: uint64(blockNum), Contract: string(contract)})
	if err != nil {
		return nil, fmt.Errorf("unable to get ABI for %q: %w", contract, err)
	}

	abi := new(zsw.ABI)
	if err = zsw.UnmarshalBinary(response.RawAbi, abi); err != nil {
		return nil, fmt.Errorf("unable to unmarshal ABI for %q: %w", contract, err)
	}

	return abi, nil
}
