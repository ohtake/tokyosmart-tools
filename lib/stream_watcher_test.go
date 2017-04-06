package lib

import (
	"testing"
)

func TestStreamWatcher_ShoudListAllFilesAtFirst(t *testing.T) {
	ch := make(chan string, 10)
	fetcher := NewDummyFetcherCount([]string{
		"12_20170402093004.ts\n"+
		"11_20170402092959.ts\n",
	})
	w := NewStreamWatcherF(NewEndpoint("v-low-tokyo1", "2865"), ch, fetcher)
	if err := w.FetchList(); err != nil {
		t.Error(err)
	}
	if len(ch) != 2 {
		t.Error(len(ch))
	}
	if f := <-ch; f != "11_20170402092959.ts" {
		t.Error(f)
	}
	if f := <-ch; f != "12_20170402093004.ts" {
		t.Error(f)
	}
}

func TestStreamWatcher_ShoudListNewFilesWhenSecondFetch(t *testing.T) {
	ch := make(chan string, 10)
	fetcher := NewDummyFetcherCount([]string{
		"12_20170402093004.ts\n"+
		"11_20170402092959.ts\n",
		"13_20170402093009.ts\n"+
		"12_20170402093004.ts\n",
	})
	w := NewStreamWatcherF(NewEndpoint("v-low-tokyo1", "2865"), ch, fetcher)
	if err := w.FetchList(); err != nil {
		t.Error(err)
	}
	if len(ch) != 2 {
		t.Error(len(ch))
	}
	<-ch
	<-ch
	if err := w.FetchList(); err != nil {
		t.Error(err)
	}
	if len(ch) != 1 {
		t.Error(len(ch))
	}
	if f := <-ch; f != "13_20170402093009.ts" {
		t.Error(f)
	}
}
