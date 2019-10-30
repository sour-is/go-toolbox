package mercury

import (
	"bufio"
	"io"
	"strings"

	"sour.is/x/toolbox/log"
)

func parseText(body io.Reader) (config SpaceMap, err error) {
	config = make(SpaceMap)

	var space string
	var name string
	var tags []string
	var notes []string
	var seq uint64

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			notes = append(notes, strings.TrimPrefix(line, "# "))
			log.Debug("NOTE ", notes)
			continue
		}

		if strings.HasPrefix(line, "@") {
			var c *Space
			var ok bool

			sp := strings.Fields(strings.TrimPrefix(line, "@"))
			log.Debug("SPACE ", sp)
			space = sp[0]

			if c, ok = config[space]; !ok {
				c = &Space{Space: space}
			}

			c.Notes = append(make([]string, 0, len(notes)), notes...)
			c.Tags = append(make([]string, 0, len(sp[1:])), sp[1:]...)

			config[space] = c
			notes = notes[:0]
			tags = tags[:0]

			continue
		}

		if space == "" {
			continue
		}

		sp := strings.SplitN(line, ":", 2)
		if len(sp) < 2 {
			continue
		}

		if strings.TrimSpace(sp[0]) == "" {
			var c *Space
			var ok bool

			if c, ok = config[space]; !ok {
				c = &Space{Space: space}
			}

			c.List[len(c.List)-1].Values = append(c.List[len(c.List)-1].Values, sp[1])
			config[space] = c

			continue
		}

		fields := strings.Fields(sp[0])
		name = fields[0]
		if len(fields) > 1 {
			tags = fields[1:]
		}

		var c *Space
		var ok bool

		if c, ok = config[space]; !ok {
			c = &Space{Space: space}
		}

		seq++
		c.List = append(
			c.List,
			Value{
				Seq:    seq,
				Name:   name,
				Tags:   append(make([]string, 0, len(tags)), tags...),
				Notes:  append(make([]string, 0, len(notes)), notes...),
				Values: []string{sp[1]},
			},
		)
		config[space] = c

		notes = notes[:0]
		tags = tags[:0]
	}

	log.Debug(config.ToArray().String())

	if err = scanner.Err(); err != nil {
		log.Error("reading standard input:", err)
	}

	return
}
