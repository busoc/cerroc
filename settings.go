package main

import (
	"log"
	"time"
)

const (
	ROCON  = "ROCON"
	ROCOFF = "ROCOFF"
	CERON  = "CERON"
	CEROFF = "CEROFF"
)

const (
	ALLIOP = "alliop.txt"
	INSTR  = "instrlist.txt"
)

var (
	ExecutionTime   time.Time
	DefaultBaseTime time.Time
)

func dumpSettings(d delta, fs fileset) {
	log.Printf("%s-%s (build: %s)", Program, Version, BuildTime)
	log.Printf("settings: AZM duration: %s", d.AZM.Duration)
	log.Printf("settings: ROCON time: %s", d.Rocon.Duration)
	log.Printf("settings: ROCOFF time: %s", d.Rocoff.Duration)
	log.Printf("settings: CER time: %s", d.Cer.Duration)
	log.Printf("settings: CERON time: %s", d.Ceron.Duration)
	log.Printf("settings: CEROFF time: %s", d.Ceroff.Duration)
	log.Printf("settings: CER crossing duration: %s", d.Intersect.Duration)
}

type Duration struct {
	time.Duration
}

func (d *Duration) String() string {
	return d.Duration.String()
}

func (d *Duration) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err == nil {
		d.Duration = v
	}
	return err
}

type delta struct {
	Rocon     Duration `toml:"rocon"`
	Rocoff    Duration `toml:"rocoff"`
	Ceron     Duration `toml:"ceron"`
	Ceroff    Duration `toml:"ceron"`
	Margin    Duration `toml:"margin"`
	Cer       Duration `toml:"cer"`
	Wait      Duration `toml:"wait"`
	Intersect Duration `toml:"crossing"`
	AZM       Duration `toml:"azm"`
	Saa       Duration `toml:"saa"`

	CerBefore    Duration `toml:"cer-before"`
	CerAfter     Duration `toml:"cer-after"`
	CerBeforeRoc Duration `toml:"cer-before-roc"`
	CerAfterRoc  Duration `toml:"cer-after-roc"`
}

type fileset struct {
	Path string `toml:"-"`

	Rocon  string `toml:"rocon"`
	Rocoff string `toml:"rocoff"`
	Ceron  string `toml:"ceron"`
	Ceroff string `toml:"ceroff"`
	Keep   bool   `toml:"-"`

	Alliop    string `toml:"-"`
	Instrlist string `toml:"-"`
}

func (f fileset) CanROC() bool {
	return f.Rocon != "" && f.Rocoff != ""
}

func (f fileset) CanCER() bool {
	return f.Ceron != "" && f.Ceroff != ""
}

func (f fileset) Empty() bool {
	return f.Rocon == "" && f.Rocoff == "" && f.Ceron == "" && f.Ceroff == ""
}

func (f fileset) Can() error {
	if (f.Rocon == "" && f.Rocoff != "") || (f.Rocon != "" && f.Rocoff == "") {
		return missingFile("ROC")
	}
	if f.Rocon == f.Rocoff {
		return sameFile("ROC")
	}
	if (f.Ceron == "" && f.Ceroff != "") || (f.Ceron != "" && f.Ceroff == "") {
		return missingFile("CER")
	}
	if f.Ceron == f.Ceroff {
		return sameFile("CER")
	}
	if f.Empty() {
		return genericErr("no command files given")
	}
	return nil
}
