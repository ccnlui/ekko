package transceiver

type delayer struct {
	interval  int64
	processed int64
	start     int64
	next      int64
}

func (d *delayer) onScheduleSend(t int64) bool {
	if d.next < t {
		d.processed++
		if d.start == 0 {
			d.start = (t / 1e9) * 1e9
		}
		d.next = d.start + int64(d.interval)*int64(d.processed)
		return true
	}
	return false
}

func (d *delayer) reset() {
	d.processed = 0
	d.start = 0
}

func (d *delayer) setInterval(interval int64) {
	if interval != d.interval {
		d.next = 0
		d.interval = interval
	}
}
