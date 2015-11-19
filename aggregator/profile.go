package aggregator

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/cover"
	"os"
)

type CoverProfile struct {
	profiles []*cover.Profile
}

func (c *CoverProfile) compareBlock(b *cover.ProfileBlock, b2 *cover.ProfileBlock) bool {
	return b.StartLine == b2.StartLine && b.EndLine == b2.EndLine &&
		b.StartCol == b2.StartCol && b.EndCol == b2.EndCol &&
		b.NumStmt == b2.NumStmt
}

func (c *CoverProfile) Aggregate(inputFile string) error {
	profiles, err := cover.ParseProfiles(inputFile)
	if err != nil {
		return err
	} else if c.profiles == nil {
		c.profiles = profiles
		return nil
	}
	for _, newProfile := range profiles {
		found := false
		for pIdx, p := range c.profiles {
			// Find the matching file
			if newProfile.FileName == p.FileName {
				found = true
				if p.Mode != newProfile.Mode {
					return fmt.Errorf("Mismatching count mode: %s vs %s", p.Mode, newProfile.Mode)
				}
				// Now merge block profiles
				for _, newBlock := range newProfile.Blocks {
					// So... find the matching block
					for bIdx, b := range p.Blocks {
						// We can't compare the blocks as a whole, since the count will differ
						if c.compareBlock(&newBlock, &b) {
							// Depending on the mode, we need to toggle or sum the count
							if p.Mode == "set" {
								if newBlock.Count > 0 {
									c.profiles[pIdx].Blocks[bIdx].Count = newBlock.Count
								}
							} else {
								c.profiles[pIdx].Blocks[bIdx].Count += newBlock.Count
							}
							// We found the correct block, go for the next new block
							break
						}
					}
				}
				// We processed the correct file, now process the next one
				break
			}
		}
		if found == false {
			c.profiles = append(c.profiles, newProfile)
		}
	}
	return nil
}

func (c *CoverProfile) Write(outputFile string) error {
	if c.profiles == nil || len(c.profiles) == 0 {
		return nil
	}
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	_, err = w.WriteString(fmt.Sprintf("mode: %s\n", c.profiles[0].Mode))
	if err != nil {
		return err
	}
	for _, profile := range c.profiles {
		for _, block := range profile.Blocks {
			buff := fmt.Sprintf("%s:%d.%d,%d.%d %d %d\n", profile.FileName, block.StartLine, block.StartCol,
				block.EndLine, block.EndCol, block.NumStmt, block.Count)
			//fmt.Println("Writing:", buff)
			_, err := w.WriteString(buff)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
