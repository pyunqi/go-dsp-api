package bid

import (
	"log"

	"github.com/mxmCherry/openrtb"
	"golang.org/x/net/context"

	"github.com/satoshi03/go-dsp-api/common/consts"
	"github.com/satoshi03/go-dsp-api/common/errors"
	"github.com/satoshi03/go-dsp-api/data"
)

func bid(ctx context.Context, br *openrtb.BidRequest) []*data.Ad {
	var selected []*data.Ad
	for _, imp := range br.Imp {
		ad, ok := getAd(ctx, &imp)
		if ok {
			ad.ImpID = imp.ID
			selected = append(selected, &ad)
		}
	}
	return selected
}

func getAd(ctx context.Context, imp *openrtb.Imp) (data.Ad, bool) {
	// Get Index having candidate ad list
	index, err := data.GetIndex(ctx, imp)
	if err != nil {
		// For debug
		log.Println(err)
		return data.Ad{}, false
	}

	// Find valid ad having max score(=ecpm) in index
	for i := range index {
		if err := validate(imp, &index[i]); err == nil {
			// Valid ad was found
			return index[i], true
		} else {
			// For debug
			log.Println(err)
		}
	}

	// No valid ad was found
	return data.Ad{}, false
}

func validate(imp *openrtb.Imp, ad *data.Ad) error {
	// Native Ad is not supported
	if imp.Native != nil {
		return errors.NoSupportError{"native"}
	}
	// Video Ad is not supported
	if imp.Video != nil {
		return errors.NoSupportError{"video"}
	}
	// Check bid currency
	if imp.BidFloorCur != "" && imp.BidFloorCur != consts.DefaultBidCur {
		return errors.InvalidCurError
	}
	// Check bid price greater than bid floor price
	if ad.CalcBidPrice() <= imp.BidFloor {
		return errors.LowPriceError
	}

	return nil
}