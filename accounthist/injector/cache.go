package injector

import (
	"github.com/zhongshuwen/histnew/accounthist"
)

func (i *Injector) UpdateSeqData(key accounthist.Facet, seqData accounthist.SequenceData) {
	i.cacheSeqData[key.String()] = seqData
}
