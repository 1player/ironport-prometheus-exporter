package main

import (
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

const SampleStatus = `
<status build="phoebe 10.0.0-203" hostname="hostname" timestamp="20170329115551">
  <birth_time timestamp="20161011101643 (169d 1h 39m 8s)"/>
  <last_counter_reset timestamp=""/>
  <system status="online" />
  <oldest_message secs="261369" mid="10836353" />
  <features>
    <feature name="McAfee" time_remaining="2592000" />
    <feature name="Sophos" time_remaining="45341935" />
    <feature name="File Analysis" time_remaining="-46298076" />
    <feature name="Bounce Verification" time_remaining="Dormant/Perpetual" />
    <feature name="IronPort Anti-Spam" time_remaining="45341935" />
    <feature name="IronPort Email Encryption" time_remaining="2592000" />
    <feature name="RSA Email Data Loss Prevention" time_remaining="2592000" />
    <feature name="File Reputation" time_remaining="-46298076" />
    <feature name="Incoming Mail Handling" time_remaining="Dormant/Perpetual" />
    <feature name="Outbreak Filters" time_remaining="45856372" />
  </features>
  <counters>
    <counter name="inj_msgs"
        reset="10741905"
        uptime="6037436"
        lifetime="10741905" />
    <counter name="inj_recips"
        reset="16324603"
        uptime="10129149"
        lifetime="16324603" />
    <counter name="gen_bounce_recips"
        reset="3771568"
        uptime="3595210"
        lifetime="3771568" />
    <counter name="rejected_recips"
        reset="130438"
        uptime="82979"
        lifetime="130438" />
    <counter name="dropped_msgs"
        reset="164522"
        uptime="67763"
        lifetime="164522" />
    <counter name="soft_bounced_evts"
        reset="41422087"
        uptime="31431593"
        lifetime="41422087" />
    <counter name="completed_recips"
        reset="16116230"
        uptime="10021778"
        lifetime="16116230" />
    <counter name="hard_bounced_recips"
        reset="7232738"
        uptime="7027240"
        lifetime="7232738" />
    <counter name="dns_hard_bounced_recips"
        reset="605549"
        uptime="597735"
        lifetime="605549" />
    <counter name="5xx_hard_bounced_recips"
        reset="5172743"
        uptime="5094589"
        lifetime="5172743" />
    <counter name="filter_hard_bounced_recips"
        reset="0"
        uptime="0"
        lifetime="0" />
    <counter name="expired_hard_bounced_recips"
        reset="1454446"
        uptime="1334916"
        lifetime="1454446" />
    <counter name="other_hard_bounced_recips"
        reset="0"
        uptime="0"
        lifetime="0" />
    <counter name="delivered_recips"
        reset="8883492"
        uptime="2994538"
        lifetime="8883492" />
    <counter name="deleted_recips"
        reset="0"
        uptime="0"
        lifetime="0" />
    <counter name="global_unsub_hits"
        reset="0"
        uptime="0"
        lifetime="0" />
  </counters>
  <current_ids
    message_id="10875286"
    injection_conn_id="17125687"
    delivery_conn_id="17502919" />
  <rates>
    <rate name="inj_msgs"
      last_1_min="601"
      last_5_min="945"
      last_15_min="1084" />
    <rate name="inj_recips"
      last_1_min="1091"
      last_5_min="1434"
      last_15_min="1409" />
    <rate name="soft_bounced_evts"
      last_1_min="94"
      last_5_min="172"
      last_15_min="132" />
    <rate name="completed_recips"
      last_1_min="1226"
      last_5_min="1410"
      last_15_min="1388" />
    <rate name="hard_bounced_recips"
      last_1_min="0"
      last_5_min="12"
      last_15_min="20" />
    <rate name="delivered_recips"
      last_1_min="1226"
      last_5_min="1398"
      last_15_min="1368" />
  </rates>
  <gauges>
    <gauge name="ram_utilization" current="6" />
    <gauge name="total_utilization" current="8" />
    <gauge name="cpu_utilization" current="0" />
    <gauge name="av_utilization" current="0" />
    <gauge name="case_utilization" current="2" />
    <gauge name="bm_utilization" current="0" />
    <gauge name="disk_utilization" current="0" />
    <gauge name="resource_conservation" current="0" />
    <gauge name="log_used" current="14" />
    <gauge name="log_available" current="388G" />
    <gauge name="conn_in" current="1" />
    <gauge name="conn_out" current="0" />
    <gauge name="active_recips" current="196" />
    <gauge name="unattempted_recips" current="70" />
    <gauge name="attempted_recips" current="126" />
    <gauge name="msgs_in_work_queue" current="0" />
    <gauge name="dests_in_memory" current="625" />
    <gauge name="kbytes_used" current="61506" />
    <gauge name="kbytes_free" current="34541502" />
    <gauge name="msgs_in_policy_virus_outbreak_quarantine" current="104" />
    <gauge name="kbytes_in_policy_virus_outbreak_quarantine" current="8390" />
    <gauge name="reporting_utilization" current="0" />
    <gauge name="quarantine_utilization" current="0" />
  </gauges>
</status>
`

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestParseStatus(c *C) {
	expectedStatus := Status{
		System: SystemStatus{
			Status: "online",
		},
		Counters: []Counter{
			Counter{Name: "inj_msgs", Lifetime: 10741905},
			Counter{Name: "inj_recips", Lifetime: 16324603},
			Counter{Name: "gen_bounce_recips", Lifetime: 3771568},
			Counter{Name: "rejected_recips", Lifetime: 130438},
			Counter{Name: "dropped_msgs", Lifetime: 164522},
			Counter{Name: "soft_bounced_evts", Lifetime: 41422087},
			Counter{Name: "completed_recips", Lifetime: 16116230},
			Counter{Name: "hard_bounced_recips", Lifetime: 7232738},
			Counter{Name: "dns_hard_bounced_recips", Lifetime: 605549},
			Counter{Name: "5xx_hard_bounced_recips", Lifetime: 5172743},
			Counter{Name: "filter_hard_bounced_recips", Lifetime: 0},
			Counter{Name: "expired_hard_bounced_recips", Lifetime: 1454446},
			Counter{Name: "other_hard_bounced_recips", Lifetime: 0},
			Counter{Name: "delivered_recips", Lifetime: 8883492},
			Counter{Name: "deleted_recips", Lifetime: 0},
			Counter{Name: "global_unsub_hits", Lifetime: 0},
		},
		Gauges: []Gauge{
			Gauge{Name: "ram_utilization", Current: 6},
			Gauge{Name: "total_utilization", Current: 8},
			Gauge{Name: "cpu_utilization", Current: 0},
			Gauge{Name: "av_utilization", Current: 0},
			Gauge{Name: "case_utilization", Current: 2},
			Gauge{Name: "bm_utilization", Current: 0},
			Gauge{Name: "disk_utilization", Current: 0},
			Gauge{Name: "resource_conservation", Current: 0},
			Gauge{Name: "log_used", Current: 14},
			Gauge{Name: "log_available", Current: 416611827712},
			Gauge{Name: "conn_in", Current: 1},
			Gauge{Name: "conn_out", Current: 0},
			Gauge{Name: "active_recips", Current: 196},
			Gauge{Name: "unattempted_recips", Current: 70},
			Gauge{Name: "attempted_recips", Current: 126},
			Gauge{Name: "msgs_in_work_queue", Current: 0},
			Gauge{Name: "dests_in_memory", Current: 625},
			Gauge{Name: "kbytes_used", Current: 61506},
			Gauge{Name: "kbytes_free", Current: 34541502},
			Gauge{Name: "msgs_in_policy_virus_outbreak_quarantine", Current: 104},
			Gauge{Name: "kbytes_in_policy_virus_outbreak_quarantine", Current: 8390},
			Gauge{Name: "reporting_utilization", Current: 0},
			Gauge{Name: "quarantine_utilization", Current: 0},
		},
	}

	status, err := parseStatus(strings.NewReader(SampleStatus))
	c.Assert(err, IsNil)
	c.Assert(*status, DeepEquals, expectedStatus)
}
