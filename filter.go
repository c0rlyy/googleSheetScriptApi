package main

import "fmt"

type Sheet struct {
	values [][]any
}

func (sh *Sheet) ConvertToMap() ([]map[any]any, error) {
	if len(sh.values) == 0 {
		return nil, fmt.Errorf("no data in sheet")
	}
	headlessValues := sh.values[1:]
	headers := sh.values[0]

	sheetContainter := make([]map[any]any, 0, len(headlessValues))
	for _, val := range headlessValues {
		rowContainer := make(map[any]any)
		for j, header := range headers {
			if j >= len(val) {
				rowContainer[header] = " "
			} else {
				rowContainer[header] = val[j]
			}
		}
		sheetContainter = append(sheetContainter, rowContainer)
	}
	return sheetContainter, nil
}

func (sh *Sheet) ConvertToValues(sheetContainer *[]map[any]any, headersOrdered []any) error {
	sheetValues := make([][]any, 0, len(*sheetContainer))
	sheetValues = append(sheetValues, headersOrdered)

	for _, sheetRow := range *sheetContainer {
		row := make([]any, 0, len(sheetRow))
		for _, val := range headersOrdered {
			cell := sheetRow[val]
			row = append(row, cell)
		}
		sheetValues = append(sheetValues, row)
	}

	sh.values = sheetValues
	return nil
}

func FilterSheetContainter(sheetContainter *[]map[any]any, config *Config, keyValue string, serialNumberVal string) (tagetValue any, err error) {
	targetValueHeader := config.TargetValueHeader
	keyHeader := config.KeyIDHeader
	serialKeyHeader := config.SerialKeyHeader

	rowKeys := (*sheetContainter)[0]
	if _, ok := rowKeys[keyHeader]; !ok {
		return nil, fmt.Errorf("key header is not present in the Sheet")
	}
	if _, ok := rowKeys[serialKeyHeader]; !ok {
		return nil, fmt.Errorf("serial number header is not present in the Sheet ")
	}
	for _, row := range *sheetContainter {
		key := row[keyHeader]
		serialNumber := row[serialKeyHeader]
		if key == keyValue && serialNumber == serialNumberVal {
			row[keyHeader] = row[targetValueHeader]
			return row[targetValueHeader], nil
		}
	}
	return nil, fmt.Errorf("provided  combinations of key and target were not found in the sheet")
}

func ReplaceSheetKey(sheetContainter *[]map[any]any, config *Config, valueToReplace any, keyValue string) (err error) {
	keyHeader := config.SecondSheetKeyIDHeader

	for _, row := range *sheetContainter {
		if row[keyHeader] == keyValue {
			row[keyHeader] = valueToReplace
			return nil
		}
	}
	return fmt.Errorf("no mathcing key was found in the sheet")
}
