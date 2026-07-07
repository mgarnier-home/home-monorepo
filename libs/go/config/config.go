package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"mgarnier11.fr/go/libs/utils"
)

func GetConfig[ConfigType any](cfg *ConfigType) []error {
	utils.InitEnvFromFile()

	configErrors := []error{}

	t := reflect.TypeOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		// On parse les fields de la struct autoscalerConfig pour récupérer les valeurs des variables d'environnement correspondantes
		field := t.Field(i)
		key := field.Tag.Get("key")
		// Si la clé est vide, on ignore ce champ
		if key == "" {
			continue
		}
		defaultValue := field.Tag.Get("default-value")
		required := field.Tag.Get("required") == "true"

		// On déclare les variables value et err en dehors du switch pour pouvoir les utiliser après le switch
		var value any
		var err error

		switch field.Type.Kind() {
		case reflect.Int:
			// Quand le type est int, on convertit la valeur par défaut en int avant de l'utiliser
			defaultInt, err := strconv.Atoi(defaultValue)
			if err != nil {
				configErrors = append(configErrors, fmt.Errorf("Invalid default value for field %s: %v", field.Name, err))
				continue
			}

			err, value = utils.GetEnvValue(key, defaultInt, required)
			if err != nil {
				configErrors = append(configErrors, err)
				continue
			}

		case reflect.String:
			err, value = utils.GetEnvValue(key, defaultValue, required)
			if err != nil {
				configErrors = append(configErrors, err)
				continue
			}
		case reflect.Bool:
			err, value = utils.GetEnvValue(key, defaultValue == "true", required)
			if err != nil {
				configErrors = append(configErrors, err)
				continue
			}
		case reflect.Slice:
			err, value = utils.GetEnvValue(key, strings.Split(defaultValue, ","), required)
			if err != nil {
				configErrors = append(configErrors, err)
				continue
			}
		default:
			configErrors = append(configErrors, fmt.Errorf("Unsupported field type for field %s", field.Name))
			continue
		}

		fieldValue := reflect.ValueOf(cfg).Elem().FieldByName(field.Name)

		if fieldValue.CanSet() {
			fieldValue.Set(reflect.ValueOf(value))
		} else {
			configErrors = append(configErrors, fmt.Errorf("Cannot set field %s", field.Name))
		}
	}

	return configErrors
}
