# Kubernetes Managed Fields Utils

This is a simple package with utils to aid the management of [managed fields](https://kubernetes.io/docs/reference/using-api/server-side-apply/) of Kubernetes objects.

This is in very embrionic status, feel free to submit PRs for improvements.

## FieldsV1ToJSONPath

It normalizes the managed fields FieldV1 Raw into a string of regular JSON Type path.

Alternatively, if passed an optional flag, it can generate regular expressions instead of paths.

It is used to have two managed fieldsV1 to be compared and matched across.

## DetectExternalFieldManager

This function detects if a managed field was overwritten by other manager.

One needs to pass the name of the manager that supposedly manage the field.

This function would detect if any field was altered by other manager and return the name of the external manager (if any other manager changed that field, before or after the original).

The use case is detecting competing two managers competing for the fields. 