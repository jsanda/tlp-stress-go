package run

import "time"

type profileRunnerCfg struct{

}

type profileRunner struct{
	Population int64
}

func createRunners() *profileRunner {
	return &profileRunner{}
}

func (p *profileRunner) Populate(rows int64, done chan struct{}) {
	defer close(done)
	for p.Population <= rows {
		p.Population++
		time.Sleep(50 * time.Nanosecond)
	}
}
