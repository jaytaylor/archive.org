package archiveorg

// TODO: integration test
// timemap, err := TimeMapFor(u)

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestParseMemento(t *testing.T) {
	lines := []string{
		`<http://www.jaytaylor.com:80/>; rel="original",`,
		`<http://web.archive.org/web/timemap/link/https://jaytaylor.com>; rel="self"; type="application/link-format"; from="Sat, 31 Mar 2001 11:48:39 GMT",`,
		`<http://web.archive.org>; rel="timegate",`,
		`<http://web.archive.org/web/20010331114839/http://www.jaytaylor.com:80/>; rel="first memento"; datetime="Sat, 31 Mar 2001 11:48:39 GMT",`,
	}

	for i, line := range lines {
		if _, err := ParseMemento(line); err != nil {
			t.Errorf("[i=%v] Error parsing entry: %s (line=%v)", i, err, line)
		}
	}
}

func TestParseTimeMap(t *testing.T) {
	r := strings.NewReader(rawTimeMap)

	timemap, err := ParseTimeMap(r)
	if err != nil {
		t.Fatalf("Error parsing TimeMap fixture: %s", err)
	}

	t.Logf("num mementos: %v", len(timemap.Mementos))
	t.Logf("timemap: %+v", timemap)

	if err := validateMemento(timemap.Original); err != nil {
		t.Errorf("Error validating Original: %v", err)
	}

	if err := validateMemento(timemap.Self); err != nil {
		t.Errorf("Error validating Self: %v", err)
	}

	if err := validateMemento(timemap.TimeGate); err != nil {
		t.Errorf("Error validating TimeGate: %v", err)
	}

	for i, memento := range timemap.Mementos {
		if err := validateMemento(&memento); err != nil {
			t.Errorf("[i=%v] Error validating Memento slice element: %v", i, err)
		}
	}

	if expected, actual := 130, len(timemap.Mementos); actual != expected {
		t.Errorf("Expected number of mementos=%v but actual=%v", expected, actual)
	}
}

func validateMemento(m *Memento) error {
	if m == nil {
		return errors.New("memento pointer is nil")
	}

	if m.Rel == "" {
		return errors.New("rel attribute is empty")
	}

	switch m.Rel {
	case "original":

	case "self":
		if m.Type == nil {
			return errors.New("type attribute is nil")
		}
		if m.From == nil {
			return errors.New("from attribute is nil")
		}

	case "timegate":

	case "memento", "first memento":
		if m.Time == nil {
			return errors.New("datetime attribute is nil")
		}

	default:
		return fmt.Errorf("no validator implemented for rel=%v", m.Rel)
	}

	return nil
}

const rawTimeMap = `
<http://www.jaytaylor.com:80/>; rel="original",
<http://web.archive.org/web/timemap/link/https://jaytaylor.com>; rel="self"; type="application/link-format"; from="Sat, 31 Mar 2001 11:48:39 GMT",
<http://web.archive.org>; rel="timegate",
<http://web.archive.org/web/20010331114839/http://www.jaytaylor.com:80/>; rel="first memento"; datetime="Sat, 31 Mar 2001 11:48:39 GMT",
<http://web.archive.org/web/20010405072207/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Thu, 05 Apr 2001 07:22:07 GMT",
<http://web.archive.org/web/20011128153904/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Wed, 28 Nov 2001 15:39:04 GMT",
<http://web.archive.org/web/20020118040828/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 18 Jan 2002 04:08:28 GMT",
<http://web.archive.org/web/20020328103634/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Thu, 28 Mar 2002 10:36:34 GMT",
<http://web.archive.org/web/20020330072512/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 30 Mar 2002 07:25:12 GMT",
<http://web.archive.org/web/20020523102917/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Thu, 23 May 2002 10:29:17 GMT",
<http://web.archive.org/web/20020524040414/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 24 May 2002 04:04:14 GMT",
<http://web.archive.org/web/20020527044428/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Mon, 27 May 2002 04:44:28 GMT",
<http://web.archive.org/web/20020720090912/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 20 Jul 2002 09:09:12 GMT",
<http://web.archive.org/web/20020802092448/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 02 Aug 2002 09:24:48 GMT",
<http://web.archive.org/web/20020922093956/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 22 Sep 2002 09:39:56 GMT",
<http://web.archive.org/web/20020925000018/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 25 Sep 2002 00:00:18 GMT",
<http://web.archive.org/web/20021123194153/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 23 Nov 2002 19:41:53 GMT",
<http://web.archive.org/web/20030210032054/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 10 Feb 2003 03:20:54 GMT",
<http://web.archive.org/web/20030217084558/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 17 Feb 2003 08:45:58 GMT",
<http://web.archive.org/web/20030419110356/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 19 Apr 2003 11:03:56 GMT",
<http://web.archive.org/web/20030608084128/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Sun, 08 Jun 2003 08:41:28 GMT",
<http://web.archive.org/web/20030723164018/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Wed, 23 Jul 2003 16:40:18 GMT",
<http://web.archive.org/web/20030802082722/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 02 Aug 2003 08:27:22 GMT",
<http://web.archive.org/web/20030802231254/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Sat, 02 Aug 2003 23:12:54 GMT",
<http://web.archive.org/web/20030929094730/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 29 Sep 2003 09:47:30 GMT",
<http://web.archive.org/web/20031024003545/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 24 Oct 2003 00:35:45 GMT",
<http://web.archive.org/web/20031121131143/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 21 Nov 2003 13:11:43 GMT",
<http://web.archive.org/web/20031212060015/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 12 Dec 2003 06:00:15 GMT",
<http://web.archive.org/web/20031231061904/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Wed, 31 Dec 2003 06:19:04 GMT",
<http://web.archive.org/web/20040327033952/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Sat, 27 Mar 2004 03:39:52 GMT",
<http://web.archive.org/web/20040414005640/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 14 Apr 2004 00:56:40 GMT",
<http://web.archive.org/web/20040418032009/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 18 Apr 2004 03:20:09 GMT",
<http://web.archive.org/web/20040525055555/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Tue, 25 May 2004 05:55:55 GMT",
<http://web.archive.org/web/20040526190155/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 26 May 2004 19:01:55 GMT",
<http://web.archive.org/web/20040609204538/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 09 Jun 2004 20:45:38 GMT",
<http://web.archive.org/web/20040611005043/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 11 Jun 2004 00:50:43 GMT",
<http://web.archive.org/web/20040613010930/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 13 Jun 2004 01:09:30 GMT",
<http://web.archive.org/web/20040824042829/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Tue, 24 Aug 2004 04:28:29 GMT",
<http://web.archive.org/web/20040829111737/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 29 Aug 2004 11:17:37 GMT",
<http://web.archive.org/web/20040901152531/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Wed, 01 Sep 2004 15:25:31 GMT",
<http://web.archive.org/web/20040924124310/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 24 Sep 2004 12:43:10 GMT",
<http://web.archive.org/web/20040925220213/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 25 Sep 2004 22:02:13 GMT",
<http://web.archive.org/web/20041203204959/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 03 Dec 2004 20:49:59 GMT",
<http://web.archive.org/web/20041212042851/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 12 Dec 2004 04:28:51 GMT",
<http://web.archive.org/web/20050125040815/http://jaytaylor.com:80/>; rel="memento"; datetime="Tue, 25 Jan 2005 04:08:15 GMT",
<http://web.archive.org/web/20050203213307/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Thu, 03 Feb 2005 21:33:07 GMT",
<http://web.archive.org/web/20050311004046/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 11 Mar 2005 00:40:46 GMT",
<http://web.archive.org/web/20050406095604/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Wed, 06 Apr 2005 09:56:04 GMT",
<http://web.archive.org/web/20050819050706/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 19 Aug 2005 05:07:06 GMT",
<http://web.archive.org/web/20050826190629/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Fri, 26 Aug 2005 19:06:29 GMT",
<http://web.archive.org/web/20071114003716/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 14 Nov 2007 00:37:16 GMT",
<http://web.archive.org/web/20100417073800/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Sat, 17 Apr 2010 07:38:00 GMT",
<http://web.archive.org/web/20100417073800/http://www.jaytaylor.com:80/>; rel="memento"; datetime="Sat, 17 Apr 2010 07:38:00 GMT",
<http://web.archive.org/web/20100529133626/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 29 May 2010 13:36:26 GMT",
<http://web.archive.org/web/20100923002742/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 23 Sep 2010 00:27:42 GMT",
<http://web.archive.org/web/20101028134350/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 28 Oct 2010 13:43:50 GMT",
<http://web.archive.org/web/20101129094913/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 29 Nov 2010 09:49:13 GMT",
<http://web.archive.org/web/20110103011751/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 03 Jan 2011 01:17:51 GMT",
<http://web.archive.org/web/20110202201644/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 02 Feb 2011 20:16:44 GMT",
<http://web.archive.org/web/20110202220545/http://jaytaylor.com/>; rel="memento"; datetime="Wed, 02 Feb 2011 22:05:45 GMT",
<http://web.archive.org/web/20110307012622/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 07 Mar 2011 01:26:22 GMT",
<http://web.archive.org/web/20110408183449/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 08 Apr 2011 18:34:49 GMT",
<http://web.archive.org/web/20110512003740/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 12 May 2011 00:37:40 GMT",
<http://web.archive.org/web/20110612073325/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 12 Jun 2011 07:33:25 GMT",
<http://web.archive.org/web/20110713151123/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 13 Jul 2011 15:11:23 GMT",
<http://web.archive.org/web/20110820160054/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 20 Aug 2011 16:00:54 GMT",
<http://web.archive.org/web/20110910023058/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 10 Sep 2011 02:30:58 GMT",
<http://web.archive.org/web/20110925034448/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 25 Sep 2011 03:44:48 GMT",
<http://web.archive.org/web/20111125202637/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 25 Nov 2011 20:26:37 GMT",
<http://web.archive.org/web/20111226124651/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 26 Dec 2011 12:46:51 GMT",
<http://web.archive.org/web/20120128020241/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 28 Jan 2012 02:02:41 GMT",
<http://web.archive.org/web/20120513041242/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 13 May 2012 04:12:42 GMT",
<http://web.archive.org/web/20120614020536/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 14 Jun 2012 02:05:36 GMT",
<http://web.archive.org/web/20120715034945/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 15 Jul 2012 03:49:45 GMT",
<http://web.archive.org/web/20121110024559/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 10 Nov 2012 02:45:59 GMT",
<http://web.archive.org/web/20130128043129/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 28 Jan 2013 04:31:29 GMT",
<http://web.archive.org/web/20130316205541/http://jaytaylor.com/>; rel="memento"; datetime="Sat, 16 Mar 2013 20:55:41 GMT",
<http://web.archive.org/web/20130405061854/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 05 Apr 2013 06:18:54 GMT",
<http://web.archive.org/web/20130502102729/http://jaytaylor.com/>; rel="memento"; datetime="Thu, 02 May 2013 10:27:29 GMT",
<http://web.archive.org/web/20130523035438/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 23 May 2013 03:54:38 GMT",
<http://web.archive.org/web/20130606034602/http://jaytaylor.com/>; rel="memento"; datetime="Thu, 06 Jun 2013 03:46:02 GMT",
<http://web.archive.org/web/20140101082753/http://jaytaylor.com/>; rel="memento"; datetime="Wed, 01 Jan 2014 08:27:53 GMT",
<http://web.archive.org/web/20140330154938/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 30 Mar 2014 15:49:38 GMT",
<http://web.archive.org/web/20140430174256/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 30 Apr 2014 17:42:56 GMT",
<http://web.archive.org/web/20140517223308/http://jaytaylor.com/>; rel="memento"; datetime="Sat, 17 May 2014 22:33:08 GMT",
<http://web.archive.org/web/20140531232744/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 31 May 2014 23:27:44 GMT",
<http://web.archive.org/web/20141012064722/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 12 Oct 2014 06:47:22 GMT",
<http://web.archive.org/web/20141112050953/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 12 Nov 2014 05:09:53 GMT",
<http://web.archive.org/web/20141217224946/http://jaytaylor.com/>; rel="memento"; datetime="Wed, 17 Dec 2014 22:49:46 GMT",
<http://web.archive.org/web/20150112121149/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 12 Jan 2015 12:11:49 GMT",
<http://web.archive.org/web/20150206012931/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 06 Feb 2015 01:29:31 GMT",
<http://web.archive.org/web/20150223120053/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 23 Feb 2015 12:00:53 GMT",
<http://web.archive.org/web/20150326042050/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 26 Mar 2015 04:20:50 GMT",
<http://web.archive.org/web/20150426005920/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 26 Apr 2015 00:59:20 GMT",
<http://web.archive.org/web/20150527024618/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 27 May 2015 02:46:18 GMT",
<http://web.archive.org/web/20150627021127/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 27 Jun 2015 02:11:27 GMT",
<http://web.archive.org/web/20150728161128/http://jaytaylor.com:80/>; rel="memento"; datetime="Tue, 28 Jul 2015 16:11:28 GMT",
<http://web.archive.org/web/20150801181332/http://jaytaylor.com/>; rel="memento"; datetime="Sat, 01 Aug 2015 18:13:32 GMT",
<http://web.archive.org/web/20150828055233/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 28 Aug 2015 05:52:33 GMT",
<http://web.archive.org/web/20150928031607/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 28 Sep 2015 03:16:07 GMT",
<http://web.archive.org/web/20151029041728/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 29 Oct 2015 04:17:28 GMT",
<http://web.archive.org/web/20151124004215/http://jaytaylor.com/>; rel="memento"; datetime="Tue, 24 Nov 2015 00:42:15 GMT",
<http://web.archive.org/web/20151130034149/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 30 Nov 2015 03:41:49 GMT",
<http://web.archive.org/web/20160109181733/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 09 Jan 2016 18:17:33 GMT",
<http://web.archive.org/web/20160110114543/http://jaytaylor.com/>; rel="memento"; datetime="Sun, 10 Jan 2016 11:45:43 GMT",
<http://web.archive.org/web/20160126013415/http://jaytaylor.com/>; rel="memento"; datetime="Tue, 26 Jan 2016 01:34:15 GMT",
<http://web.archive.org/web/20160205033340/http://jaytaylor.com/>; rel="memento"; datetime="Fri, 05 Feb 2016 03:33:40 GMT",
<http://web.archive.org/web/20160209104907/http://jaytaylor.com:80/>; rel="memento"; datetime="Tue, 09 Feb 2016 10:49:07 GMT",
<http://web.archive.org/web/20160311212911/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 11 Mar 2016 21:29:11 GMT",
<http://web.archive.org/web/20160312102435/http://jaytaylor.com/>; rel="memento"; datetime="Sat, 12 Mar 2016 10:24:35 GMT",
<http://web.archive.org/web/20160318174443/http://jaytaylor.com/>; rel="memento"; datetime="Fri, 18 Mar 2016 17:44:43 GMT",
<http://web.archive.org/web/20160331222402/http://jaytaylor.com/>; rel="memento"; datetime="Thu, 31 Mar 2016 22:24:02 GMT",
<http://web.archive.org/web/20160412081000/http://jaytaylor.com:80/>; rel="memento"; datetime="Tue, 12 Apr 2016 08:10:00 GMT",
<http://web.archive.org/web/20160513072538/http://jaytaylor.com:80/>; rel="memento"; datetime="Fri, 13 May 2016 07:25:38 GMT",
<http://web.archive.org/web/20160613145620/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 13 Jun 2016 14:56:20 GMT",
<http://web.archive.org/web/20160821143108/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 21 Aug 2016 14:31:08 GMT",
<http://web.archive.org/web/20160922014524/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 22 Sep 2016 01:45:24 GMT",
<http://web.archive.org/web/20161011051941/http://jaytaylor.com/>; rel="memento"; datetime="Tue, 11 Oct 2016 05:19:41 GMT",
<http://web.archive.org/web/20161023072003/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 23 Oct 2016 07:20:03 GMT",
<http://web.archive.org/web/20161216231030/http://jaytaylor.com/>; rel="memento"; datetime="Fri, 16 Dec 2016 23:10:30 GMT",
<http://web.archive.org/web/20170415041218/http://jaytaylor.com:80/>; rel="memento"; datetime="Sat, 15 Apr 2017 04:12:18 GMT",
<http://web.archive.org/web/20170516034253/http://jaytaylor.com:80/>; rel="memento"; datetime="Tue, 16 May 2017 03:42:53 GMT",
<http://web.archive.org/web/20170607222123/http://jaytaylor.com:80/>; rel="memento"; datetime="Wed, 07 Jun 2017 22:21:23 GMT",
<http://web.archive.org/web/20170626170713/http://jaytaylor.com/>; rel="memento"; datetime="Mon, 26 Jun 2017 17:07:13 GMT",
<http://web.archive.org/web/20170807081346/http://jaytaylor.com:80/>; rel="memento"; datetime="Mon, 07 Aug 2017 08:13:46 GMT",
<http://web.archive.org/web/20170907045442/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 07 Sep 2017 04:54:42 GMT",
<http://web.archive.org/web/20170927074436/http://jaytaylor.com/>; rel="memento"; datetime="Wed, 27 Sep 2017 07:44:36 GMT",
<http://web.archive.org/web/20171026043313/http://jaytaylor.com:80/>; rel="memento"; datetime="Thu, 26 Oct 2017 04:33:13 GMT",
<http://web.archive.org/web/20171126040432/http://jaytaylor.com:80/>; rel="memento"; datetime="Sun, 26 Nov 2017 04:04:32 GMT",
<http://web.archive.org/web/20180305214959/https://jaytaylor.com/>; rel="memento"; datetime="Mon, 05 Mar 2018 21:49:59 GMT",
<http://web.archive.org/web/20180307055959/https://jaytaylor.com/>; rel="memento"; datetime="Wed, 07 Mar 2018 05:59:59 GMT",
<http://web.archive.org/web/20180329054824/http://jaytaylor.com/>; rel="memento"; datetime="Thu, 29 Mar 2018 05:48:24 GMT",
<http://web.archive.org/web/20180519054157/https://jaytaylor.com/>; rel="memento"; datetime="Sat, 19 May 2018 05:41:57 GMT",
`
