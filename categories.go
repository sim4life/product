package product

import (
	"bufio"
	"fmt"
	"reflect"
	// "sort"
	"strconv"
	"strings"
)

type Category struct {
	ID                 int    `json:"-"`
	Category_ID        int64  `json:"categoryId"`
	Parent_Category_ID int64  `json:"-"`
	Category_Name      string `json:"name"`
	Search_Synonyms    string `json:"searchSynonym"`
	Sequence           string `json:"-"`

	parent   *Category   `json:"-"`
	children []*Category `json:"-"`
}

func PrintCategoryMap(cat_map map[int64]*Category) {
	// var m map[int]string
	var keys []int64
	for k, v := range cat_map {
		keys = append(keys, k)
		fmt.Printf("key: %d and val:%+v\n", k, v)
		fmt.Printf("parent is: %+v\n", v.parent)
		if len(v.children) > 0 {
			fmt.Printf("child[0] is: %+v\n\n", v.children[0])
		}
	}
	fmt.Printf("\nGrandaddy is:%+v\n", cat_map[int64(1)])

}

func populateCategoryTree(cat_map map[int64]*Category) {
	counter := 0
	//finding the root of the tree
	var root_groc_cat *Category
	for _, v := range cat_map {
		if v.Parent_Category_ID == int64(0) {
			root_groc_cat = v
			counter++
		}
	}

	// TODO: if counter is greater than 1 then
	// throw error/panic
	fmt.Printf("root is: %+v", root_groc_cat)
	fmt.Println("found roots: ", counter)

	for _, v := range cat_map {
		parent := cat_map[v.Parent_Category_ID]
		if parent != nil {
			// parent.children = appendPtr(parent.children, v)
			parent.children = append(parent.children, v)
			v.parent = parent
			/*
				if *parent.children == nil {
					// parent.children = *[]Category{*v}
					parent.children = v
					v.parent = parent
				} else {
					*parent.children = append(*parent.children, *v)
					v.parent = parent
				}*/
		}
	}
}

func ParseCategoriesFile(scanner *bufio.Scanner) interface{} {
	fmt.Println("fun fun")

	counter := 1
	cat_id := int64(0)
	category := &Category{ID: counter}
	cat_map := make(map[int64]*Category)

	for scanner.Scan() {
		str := strings.Split(scanner.Text(), "|")
		if str[0] != "" {
			fmt.Printf("str[0]:%s ", str[0])
			if str[0] != "EOR" {
				fmt.Printf("str[1]:%s\n", str[1])

				field_kind := reflect.ValueOf(*category).FieldByName(str[0]).Kind()
				if field_kind == reflect.String {
					reflect.ValueOf(category).Elem().FieldByName(str[0]).SetString(str[1])
				} else { //field_kind == reflect.Int
					if value, err := strconv.Atoi(str[1]); err == nil {
						if str[0] == "Category_ID" {
							cat_id = int64(value)
						}
						reflect.ValueOf(category).Elem().FieldByName(str[0]).SetInt(int64(value))
					}
				}

			} else {
				// if counter%100 == 0 {
				fmt.Println("{\"index\" : {\"_id\" :", counter, ", \"_index\":\"policereport\", \"_type\":\"crimestat\"}}")
				fmt.Println("%+v\nFinish----------", category)
				// }

				cat_map[cat_id] = category
				counter++
				category = &Category{ID: counter}

				fmt.Println("%+v----------Init", category)

			}
		}
	}

	populateCategoryTree(cat_map)
	PrintCategoryMap(cat_map)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}

	return cat_map

}
