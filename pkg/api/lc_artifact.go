/*
 * Copyright (c) 2018-2020 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	immuschema "github.com/codenotary/immudb/pkg/api/schema"
	"github.com/vchain-us/ledger-compliance-go/schema"
	"github.com/vchain-us/vcn/pkg/meta"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (a Artifact) toLcArtifact() *LcArtifact {
	aR := &LcArtifact{
		// root fields
		Kind:        a.Kind,
		Name:        a.Name,
		Hash:        a.Hash,
		Size:        a.Size,
		ContentType: a.ContentType,

		// custom metadata
		Metadata: a.Metadata,
	}

	return aR
}
func ItemToLcArtifact(item *schema.ItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(item.Item.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(item.Timestamp.GetSeconds(), int64(item.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if item.ApikeyRevoked != nil {
		if item.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(item.ApikeyRevoked.GetSeconds(), int64(item.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = item.LedgerName
	return &lca, nil
}

func ZItemToLcArtifact(ie *schema.ZItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(ie.Item.Entry.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(ie.Timestamp.GetSeconds(), int64(ie.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if ie.ApikeyRevoked != nil {
		if ie.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(ie.ApikeyRevoked.GetSeconds(), int64(ie.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = ie.LedgerName
	return &lca, nil
}

func VerifiableItemExtToLcArtifact(item *schema.VerifiableItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(item.Item.Entry.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(item.Timestamp.GetSeconds(), int64(item.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if item.ApikeyRevoked != nil {
		if item.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(item.ApikeyRevoked.GetSeconds(), int64(item.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = item.LedgerName
	return &lca, nil
}

type LcArtifact struct {
	// root fields
	Uid         string    `json:"uid" yaml:"uid" vcn:"UID"`
	Kind        string    `json:"kind" yaml:"kind" vcn:"Kind"`
	Name        string    `json:"name" yaml:"name" vcn:"Name"`
	Hash        string    `json:"hash" yaml:"hash" vcn:"Hash"`
	Size        uint64    `json:"size" yaml:"size" vcn:"Size"`
	Timestamp   time.Time `json:"timestamp,omitempty" yaml:"timestamp" vcn:"Timestamp"`
	ContentType string    `json:"contentType" yaml:"contentType" vcn:"ContentType"`

	// custom metadata
	Metadata    Metadata     `json:"metadata" yaml:"metadata" vcn:"Metadata"`
	Attachments []Attachment `json:"attachments" yaml:"attachments" vcn:"Attachments"`

	Signer  string      `json:"signer" yaml:"signer" vcn:"SignerID"`
	Revoked *time.Time  `json:"revoked,omitempty" yaml:"revoked" vcn:"Apikey revoked"`
	Status  meta.Status `json:"status" yaml:"status" vcn:"Status"`
	Ledger  string      `json:"ledger,omitempty" yaml:"ledger"`
}

func (u LcUser) createArtifact(artifact Artifact, status meta.Status, attach []string) (bool, uint64, error) {

	aR := artifact.toLcArtifact()
	aR.Status = status

	aR.Signer = GetSignerIDByApiKey(u.Client.ApiKey)

	// vcn.myApiKey.{artifact hash}
	// attachment key need to have "vcn." prefix because it's handled inside cnil frontend. (attachment is listed in the UI).
	key := AppendPrefix(meta.VcnPrefix, []byte(aR.Signer))
	key = AppendSignerId(artifact.Hash, key)

	// Attachments handler
	// attachments info generation and multi kv preparation
	var aKVs []*immuschema.KeyValue
	var aRattachment []Attachment

	// map to save all the attachments with a specific label
	labelMap := make(map[string][]Attachment)
	for _, al := range attach {
		// attachment can be --attach=vscanner.result:jobid123. jobid123 is the label
		alSlice := strings.SplitN(al, ":", 2)
		a := alSlice[0]
		/** friendly label **/
		label := ""
		if len(alSlice) > 1 {
			label = alSlice[1]
		}

		// attachment
		f, err := os.Open(a)
		if err != nil {
			return false, 0, err
		}
		defer f.Close()

		fc, err := ioutil.ReadFile(a)
		if err != nil {
			return false, 0, err
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return false, 0, err
		}
		checksum := h.Sum(nil)
		hash := hex.EncodeToString(checksum)
		akey := AppendAttachment(hash, key)

		kv := &immuschema.KeyValue{
			Key:   []byte(akey),
			Value: fc,
		}

		aKVs = append(aKVs, kv)

		mime := http.DetectContentType(fc)
		at := Attachment{
			Filename: path.Base(a),
			Hash:     hash,
			Mime:     mime,
		}

		/** friendly label **/
		/* _ITEM.ATTACH.LABEL.myApiKey.{arifact hash}.vscanner.result:jobid123 */
		if label != "" {
			labelKey := meta.VcnAttachmentLabelPrefix + "." + aR.Signer + "." + artifact.Hash + "." + al
			// here is used an array to be downloadable by the same code in the attachments map use case
			attachs := []Attachment{at}
			attachmentsListJson, err := json.Marshal(attachs)
			if err != nil {
				return false, 0, err
			}
			labelKV := &immuschema.KeyValue{
				Key:   []byte(labelKey),
				Value: attachmentsListJson,
			}
			aKVs = append(aKVs, labelKV)

			// label map
			// append the attachment key in the labelMap at specific label key
			labelMap[label] = append(labelMap[label], at)
		}

		aRattachment = append(aRattachment, at)
	}

	aR.Attachments = aRattachment
	arJson, err := json.Marshal(aR)
	if err != nil {
		return false, 0, err
	}

	md := metadata.Pairs(
		meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue,
		meta.VcnLCCmdHeaderName, meta.VcnLCNotarizeCmdHeaderValue,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var txMeta *immuschema.TxMetadata
	eor := &immuschema.SetRequest{
		KVs: []*immuschema.KeyValue{
			{
				Key:   key,
				Value: arJson,
			},
		},
	}
	if len(aKVs) > 0 {
		eor.KVs = append(eor.KVs, aKVs...)
	}

	// here is built a key to retrieve in a single call all the attachment with a specific label. The value is a list of attachment keys joined by ":" separator
	for label, attachments := range labelMap {
		/* _ITEM.ATTACH.LABEL.myApiKey.{arifact hash}.jobid123 */
		labelMapKey := meta.VcnAttachmentLabelPrefix + "." + aR.Signer + "." + artifact.Hash + "." + label

		attachmentsListJson, err := json.Marshal(attachments)
		if err != nil {
			return false, 0, err
		}
		labelMapKV := &immuschema.KeyValue{
			Key:   []byte(labelMapKey),
			Value: attachmentsListJson, // attachmentKeys
		}

		eor.KVs = append(eor.KVs, labelMapKV)
	}
	txMeta, err = u.Client.SetAll(ctx, eor)
	if err != nil {
		return false, 0, err
	}
	return true, txMeta.Id, nil
}

// LoadArtifact fetches and returns an *lcArtifact for the given hash and current u, if any.
func (u *LcUser) LoadArtifact(
	hash, signerID string,
	uid string,
	tx uint64,
	gRPCMetadata map[string][]string,
) (lc *LcArtifact, verified bool, err error) {

	md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
	if len(gRPCMetadata) > 0 {
		md = metadata.Join(md, gRPCMetadata)
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if signerID == "" {
		signerID = GetSignerIDByApiKey(u.Client.ApiKey)
	}

	key := AppendPrefix(meta.VcnPrefix, []byte(signerID))
	key = AppendSignerId(hash, key)

	var jsonAr *schema.VerifiableItemExt
	if uid != "" {
		score, err := strconv.ParseFloat(uid, 64)
		if err != nil {
			return nil, false, err
		}
		zitems, err := u.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
			Set:       key,
			SeekScore: math.MaxFloat64,
			SeekAtTx:  tx,
			Limit:     1,
			MinScore:  &immuschema.Score{Score: score},
			MaxScore:  &immuschema.Score{Score: score},
			SinceTx:   math.MaxUint64,
			NoWait:    true,
		})
		if err != nil {
			return nil, false, err
		}
		if len(zitems.Items) > 0 {
			jsonAr, err = u.Client.VerifiedGetExtAt(ctx, zitems.Items[0].Item.Key, zitems.Items[0].Item.AtTx)
		} else {
			return nil, false, ErrNotFound
		}
	} else {
		jsonAr, err = u.Client.VerifiedGetExtAt(ctx, key, tx)
	}
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Message() == "data is corrupted" {
			return nil, false, ErrNotVerified
		}
		if err.Error() == "data is corrupted" {
			return nil, false, ErrNotVerified
		}
		if ok && s.Message() == "key not found" {
			return nil, false, ErrNotFound
		}
		return nil, true, err
	}

	lcArtifact, err := VerifiableItemExtToLcArtifact(jsonAr)
	if err != nil {
		return nil, false, err
	}

	return lcArtifact, true, nil
}

// GetArtifactAttachmentListByLabel returns the attachment list of an artifact and the most recent uid by a provided label and signerID
// When there are multiple attachments with same file name it adds an enumerator postfix.
func (u *LcUser) GetArtifactAttachmentListByLabel(hash string, signerID, label string) ([]Attachment, string, error) {
	if label == "" {
		return nil, "", errors.New("no attachment provided")
	}
	if hash == "" {
		return nil, "", errors.New("no artifact provided")
	}
	var attachmentList []Attachment
	var uid string
	attachmentMap, err := u.fetchAttachmentMapByLabel(hash, signerID, label)
	if err != nil {
		return nil, "", err
	}
	// map order is not guaranted so here obtain a sorted string array
	var attachDriver []string
	for k, _ := range attachmentMap {
		attachDriver = append(attachDriver, k)
	}
	sort.Strings(attachDriver)
	// reverse the driver
	last := len(attachDriver) - 1
	for i := 0; i < len(attachDriver)/2; i++ {
		attachDriver[i], attachDriver[last-i] = attachDriver[last-i], attachDriver[i]
	}
	// attachmentFileNameMap is used internally to produce a map to handle attachments with same name
	attachmentFileNameMap := make(map[string][]*Attachment)

	for _, k := range attachDriver {
		attachMapEntry := attachmentMap[k]
		// latest uid, needed to authenticate the latest notarized artifact
		if uid == "" {
			uid = k
		}
		for _, att := range attachMapEntry {
			fn := att.Filename
			if _, ok := attachmentFileNameMap[fn]; ok && len(attachmentFileNameMap[fn]) > 0 {
				// if there is a newer filename here a postfix is added. ~1,~2 ... ~N
				att.Filename = fn + "~" + strconv.Itoa(len(attachmentFileNameMap[fn]))
			}
			attachmentFileNameMap[fn] = append(attachmentFileNameMap[fn], att)
			// attachmentList contains all attachments with latest first order
			attachmentList = append(attachmentList, *att)
		}
	}
	return attachmentList, uid, nil
}

func (u *LcUser) fetchAttachmentMapByLabel(hash, signerID string, attach string) (map[string][]*Attachment, error) {

	md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if signerID == "" {
		signerID = GetSignerIDByApiKey(u.Client.ApiKey)
	}

	key := meta.VcnAttachmentLabelPrefix + "." + signerID + "." + hash + "." + attach

	/* _ITEM.ATTACH.LABEL.myApiKey.{arifact hash}.vscanner.result:jobid123 */
	/* _ITEM.ATTACH.LABEL.myApiKey.{arifact hash}.jobid123 */
	sr := &immuschema.ScanRequest{
		Prefix:  []byte(key),
		SinceTx: math.MaxUint64,
		NoWait:  true,
		Desc:    true,
	}

	res, err := u.Client.Scan(ctx, sr)
	if err != nil {
		return nil, err
	}
	if len(res.Entries) < 1 {
		return nil, errors.New("provided label does not contains entries")
	}

	attachMap := make(map[string][]*Attachment)

	for _, entry := range res.Entries {
		// ori reg ex _ITEM\.ATTACH\.LABEL\.[^.]+\.[^.]+\.(\S+:\S[^.]+|\S+)\.([0-9]+)
		var regex = regexp.MustCompile("_ITEM\\.ATTACH\\.LABEL\\.[^.]+\\.[^.]+\\.(\\S+:\\S[^.]+|\\S+)\\.([0-9]+)")
		keyAndUid := regex.FindStringSubmatch(string(entry.Key))

		if len(keyAndUid) != 3 {
			return nil, errors.New("not consistent data when retrieving uid from attachment label entry")
		}

		attachmentList := make([]*Attachment, 0)
		err = json.Unmarshal(entry.Value, &attachmentList)
		if err != nil {
			return nil, err
		}
		attachMap[keyAndUid[2]] = attachmentList
	}

	return attachMap, nil
}

func AppendPrefix(prefix string, key []byte) []byte {
	var prefixed = make([]byte, len(prefix)+1+len(key))
	copy(prefixed[0:], prefix+".")
	copy(prefixed[len(prefix)+1:], key)
	return prefixed
}

func AppendSignerId(signerId string, k []byte) []byte {
	var prefixed = make([]byte, len(k)+len(signerId)+1)
	copy(prefixed[0:], k)
	copy(prefixed[len(k):], "."+signerId)
	return prefixed
}

func AppendAttachment(attachHash string, key []byte) []byte {
	//vcn.$AssetHash.Attachment.$AttachmentHash
	var prefixed = make([]byte, len(attachHash)+len(meta.AttachmentSeparator)+len(key))
	copy(prefixed[0:], key)
	copy(prefixed[len(key):], meta.AttachmentSeparator+attachHash)
	return prefixed
}

func AppendLabel(label string, key []byte) []byte {
	//vcn.$AssetHash.Attachment.$AttachmentHash
	var prefixed = make([]byte, len(label)+len(meta.AttachmentSeparator)+len(key))
	copy(prefixed[0:], key)
	copy(prefixed[len(key):], meta.AttachmentSeparator+label)
	return prefixed
}

// DownloadAttachment download locally all the attachments linked to the assets
func (u *LcUser) DownloadAttachment(attach *Attachment, ar *LcArtifact, tx uint64, lcAttachForce bool) (err error) {

	md := metadata.Pairs(meta.VcnLCPluginTypeHeaderName, meta.VcnLCPluginTypeHeaderValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	key := AppendPrefix(meta.VcnPrefix, []byte(ar.Signer))
	key = AppendSignerId(ar.Hash, key)
	attachmentKey := AppendAttachment(attach.Hash, key)

	attachEntry, err := u.Client.VerifiedGetAt(ctx, attachmentKey, tx)
	if err != nil {
		return err
	}
	if _, err := os.Stat(attach.Filename); os.IsNotExist(err) || lcAttachForce {
		return ioutil.WriteFile(attach.Filename, attachEntry.Value, 0644)
	}
	return fmt.Errorf("attachment %s already present on disk. Use --force to overwrite silently", attach.Filename)
}

// Date returns a RFC3339 formatted string of verification time (v.Timestamp), if any, otherwise an empty string.
func (lca *LcArtifact) Date() string {
	if lca != nil {
		ut := lca.Timestamp.UTC()
		if ut.Unix() > 0 {
			return ut.Format(time.RFC3339)
		}
	}
	return ""
}
