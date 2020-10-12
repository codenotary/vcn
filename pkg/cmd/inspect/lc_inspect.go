/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package inspect

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/vchain-us/vcn/pkg/meta"

	"fmt"
	"github.com/vchain-us/vcn/pkg/cmd/internal/cli"
	"google.golang.org/grpc/metadata"

	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
)

func lcInspect(hash string, signerID string, u *api.LcUser, output string) (err error) {
	md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var items *schema.StructuredItemList

	hasher := sha256.New()
	hasher.Write([]byte(u.LcApiKey()))
	contextSignerID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if signerID == "" {
		if output == "" {
			fmt.Printf("no signer ID provided. Only the last item for each signer ID is returned\n")
		}
		items, err = u.Client.ZScan(ctx, []byte(hash))
		if err != nil {
			return err
		}
	} else {
		contextSignerID = signerID
		items, err = u.Client.History(ctx, []byte(hash))
		if err != nil {
			return err
		}
	}

	// todo @michele we should improve zscan adding functionalities to secondary indexes to retrieve the full history of a key, not only the last one
	var list = make(map[uint64]*schema.StructuredItem)
	for _, v := range items.Items {
		list[v.Index] = v
	}

	l := len(list)
	if output == "" {
		fmt.Printf(
			`current signer ID "%s"
%d notarizations found for "%s"

`,
			contextSignerID, l, hash)
	}

	results := make([]types.LcResult, l)
	var i = 0
	for _, v := range list {
		var lca api.LcArtifact
		err = json.Unmarshal(v.Value.Payload, &lca)
		if err != nil {
			return err
		}
		results[i] = *types.NewLcResult(&lca, true)
		if err != nil {
			results[i].AddError(err)
		}
		i++
	}

	return cli.PrintLcSlice(output, results)
}
