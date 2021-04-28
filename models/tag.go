package models

import "fmt"

type Tag struct {
	ID		int
	Label 	string
}

var (
	tags 		[]*Tag
	nextTagID 	= 1
)

func GetTags() []*Tag{
	return tags
}

//Return the tag in case it exists
func GetTagByLabel(l string) (Tag, error) {
	b,i,e := ExistTagByLabel(l)
	if b {
		return *tags[i], nil
	}

	return Tag{}, e
}

//Add to the list and return the added Tag
func AddTag(t Tag) (Tag, error) {
	//Validation
	if t.ID != 0 {
		return Tag{}, fmt.Errorf("Tag must not contain ID")
	}
	if t.Label == "" {
		return Tag{}, fmt.Errorf("Tag label must not be empty")
	}

	//Test if tag already exists
	b,i,e := ExistTagByLabel(t.Label)
	if b {
		return *tags[i], e
	}

	//Add new tag
	t.ID = nextTagID
	nextTagID++
	tags = append(tags, &t)
	return t,nil
}

//Return if a tag exists on the list
func ExistTagByLabel(l string) (bool, int, error) {
	for i,t := range tags {
		if t.Label == l{
			return true,i,nil
		}
	}
	return false, -1, fmt.Errorf("Tag '%v' not found", l)
}