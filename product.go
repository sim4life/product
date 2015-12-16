package product

import (
	"bufio"
	// "encoding/json"
	//"flag"
	"fmt"
	//elastigo "github.com/mattbaird/elastigo/lib"
	//"log"
	// "os"
	"reflect"
	"strconv"
	"strings"
)

type PrintStructVals interface {
	PrintValues()
}

type Products []Product

func (p Products) PrintValues() {
	for ind, product := range p {
		fmt.Printf("product[%d]:\n%+v\n", ind, product)
	}
}

type Product struct {
	ID                        int        `json:"-"`
	Product_ID                int64      `json:"productId"`
	Name                      string     `json:"name"`
	Item_Code                 string     `json:"itemCode"`
	Description               string     `json:"description"`
	Record_State              int64      `json:"-"`
	Used_In_Recipe_ID         []int64    `json:"-"`
	Dep_Comm                  string     `json:"-"`
	Groceries_Category_ID     []int64    `json:"-"`
	Groceries_Categories      [][]string `json:"groceriesCategories"`
	Groceries_Search_Synonyms []string   `json:"gorceriesSearchSynonyms"`
	Ideas_Category_ID         []int64    `json:"-"`
	Ideas_Categories          [][]string `json:"ideasCategories"`
	Ideas_Search_Synonyms     []string   `json:"ideasSearchSynonyms"`
}

type SliceField struct {
	field  string
	values []int64
	isOld  bool
}

/*
var (
	host *string = flag.String("host", "localhost", "Elasticsearch Host")
)
*/
/*
func readStringField(product *Product, field, value string) string {
	fieldName := reflect.ValueOf(product).Elem().FieldByName(field)
	if fieldName.IsValid() {
		fieldName.SetString(value)
		return ""
	} else {
		// unknowns = append(unknowns, field)
		return field
	}
}
*/

func setCurrentSliceField(product *Product, sliceField *SliceField) {
	fmt.Printf("-/-/--sliceField.isOld is:%v\n", sliceField.isOld)
	if sliceField.isOld {
		fmt.Printf("str setting values:%v\n", sliceField.values)
		reflect.ValueOf(product).Elem().FieldByName(sliceField.field).Set(reflect.ValueOf(sliceField.values))
		sliceField.isOld = false
	}

	//return sliceField

}

func readSliceField(sliceField *SliceField, product *Product, field string, value int64) SliceField {
	fmt.Printf("sField:%s while field:%s == %v\n", sliceField.field, field, sliceField.isOld)
	fmt.Printf("the condition is: %v\n", (sliceField.field != field || !sliceField.isOld))

	// if !sliceField.isOld {
	if sliceField.field != field || !sliceField.isOld {
		setCurrentSliceField(product, sliceField)

		// sliceField.values = []int64{value}
		// sliceField.field = field
		sliceField = &SliceField{field, []int64{value}, true}
		fmt.Printf("  --Inner sliceField is: %+v\n", sliceField)
	} else {
		sliceField.values = append(sliceField.values, value)
	}

	fmt.Printf(" -----sliceField is: %+v\n", sliceField)
	return *sliceField
}

func readSetField(product *Product, field, value string, setField func(field reflect.Value, value string)) string {
	fieldName := reflect.ValueOf(product).Elem().FieldByName(field)
	if fieldName.IsValid() {
		// fieldName.SetInt(int64(value))
		setField(fieldName, value)
		return ""
	} else {
		return field
	}

}

func setStrIntField(sliceField *SliceField, unknowns []string, product *Product, field, value string,
	setField func(field reflect.Value, value string)) {

	setCurrentSliceField(product, sliceField)
	/*
		if sliceField.isOld {
		// if sliceField.field != field {
			// setCurrentSliceField(*product, *sliceField)
			fmt.Printf("str setting values:%v\n", sliceField.values)
			reflect.ValueOf(product).Elem().FieldByName(sliceField.field).Set(reflect.ValueOf(sliceField.values))
			sliceField.isOld = false
		}
	*/
	unknown := readSetField(product, field, value, setField)
	if unknown != "" {
		unknowns = append(unknowns, unknown)
	}

}

func ParseProductFile(scanner *bufio.Scanner) PrintStructVals {
	counter := 1
	product := Product{ID: counter}
	// products := make([]Product, 0)
	var products Products
	var sliceField SliceField
	// old_field := ""
	unknowns := make([]string, 0)
	// unknown := ""
	// old_field_cnt := 0
	// values := make([]int64, 1)
	for scanner.Scan() {
		str := strings.Split(scanner.Text(), "|")
		// fmt.Printf("str[0]:%s ", str[0])
		if str[0] != "EOR" {
			// fmt.Printf("str[1]:%s\n", str[1])
			// fmt.Printf("\nstr[0]:%s | str[1]:%s", str[0], str[1])
			// typeOf := reflect.TypeOf(&product).Elem().FieldByName(str[0])
			/*
				kindOf := reflect.ValueOf(product).FieldByName(str[0]).Kind()
				valueOf := reflect.ValueOf(product).FieldByName(str[0]).Type()
				valueOfElem := reflect.ValueOf(&product).Elem().FieldByName(str[0])

				fmt.Printf("\nvalueOf is:%+v", valueOf)
				fmt.Printf(" typeOf is:%s", valueOfElem.Type())
				fmt.Printf(" kindOf is:%s", kindOf)
				fmt.Printf(" valueOfElem is:%+v\n", valueOfElem)

				// fmt.Printf("\nvalueOf == int is:%+b", valueOf == int)
				// fmt.Printf("\nvalueOf == int.type is:%+b", valueOf == int.type())
				// fmt.Printf(" typeOf is:%s", valueOfElem.Type())
				fmt.Printf(" kindOf == reflect.Slice is:%b", kindOf == reflect.Slice)
				// fmt.Printf(" valueOfElem is:%+v\n", valueOfElem)
			*/
			// value, err := strconv.Atoi(str[1])
			// if err != nil {
			field_kind := reflect.ValueOf(product).FieldByName(str[0]).Kind()
			// fmt.Printf("  --==++---field_kind is:%v\n", field_kind)
			if field_kind == reflect.String || field_kind == reflect.Invalid {
				setStrIntField(&sliceField, unknowns, &product, str[0], str[1],
					func(field reflect.Value, value string) {
						field.SetString(value)
					})
				/*
					if old_field != "" {
						fmt.Printf("str setting values:%v\n", values)
						reflect.ValueOf(&product).Elem().FieldByName(old_field).Set(reflect.ValueOf(values))
						old_field = ""
					}
						unknown = readSetField(&product, str[0], str[1],
							func(field reflect.Value, value string) {
								field.SetString(value)
							})
						if unknown != "" {
							unknowns = append(unknowns, unknown)
						}

				*/
			} else {
				if value, err := strconv.ParseInt(str[1], 10, 64); err == nil {
					//TODO: get the value of field from the struct and then manipulate it
					// instead of relying on data input as contiguous
					if field_kind == reflect.Slice {
						sliceField = readSliceField(&sliceField, &product, str[0], value)
						// fmt.Printf("values len is:%d", len(values))
					} else {
						setStrIntField(&sliceField, unknowns, &product, str[0], str[1],
							func(field reflect.Value, value string) {
								int_val, _ := strconv.ParseInt(value, 10, 64)
								field.SetInt(int_val)
							})
						/*
							if old_field != "" {
								fmt.Printf("int setting values:%v\n", values)
								reflect.ValueOf(&product).Elem().FieldByName(old_field).Set(reflect.ValueOf(values))
								old_field = ""
							}
							fmt.Printf(" value: %d ", value)
							unknown := readSetField(&product, str[0], str[1],
								func(field reflect.Value, value string) {
									int_val, _ := strconv.ParseInt(value, 10, 64)
									field.SetInt(int_val)
								})
							if unknown != "" {
								unknowns = append(unknowns, unknown)
							}
						*/
					}
				}
			}
			/* else {
				fieldName := reflect.ValueOf(&product).Elem().FieldByName(str[0])
				if !fieldName.IsValid() {
					unknowns = append(unknowns, str[0])
					setCurrentSliceField(&product, &sliceField)
				}
			}*/

			// b, _ := json.Marshal(product)
		} else {
			//flush any pending slice data
			setCurrentSliceField(&product, &sliceField)
			/*
				if counter%20 == 0 {
					fmt.Println("\n\n{\"index\" : {\"_id\" :", counter, ", \"_index\":\"policereport\", \"_type\":\"crimestat\"}}")
					fmt.Println("%+v\nFinish----------", product)
				}
			*/
			products = append(products, product)
			counter++
			product = Product{ID: counter}
			fmt.Println("%+v----------Init", product)
		}
	}

	fmt.Printf("\nSuccessfully read products: %v\n", len(products))
	fmt.Printf("\nNo. of unkowns de-normalized fields are: %v\n", len(unknowns))
	/*
		for _, p := range products {
			fmt.Printf("products: %+v\n", p)
		}
	*/

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil //, err
	}
	return products

}