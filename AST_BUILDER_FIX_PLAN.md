# AST Builder Fix Plan

## Current Status Analysis

### ✅ Fixed Issues:
1. **Syntax errors** in `http.Error` calls (extra parentheses)
2. **Handler structure** (validator and create handler fields)
3. **Basic parameter parsing** (cookies, headers, request body)
4. **Route generation** (fixed path pattern)
5. **Apimodels structure** (added RequestCookies and Cookies field)
6. **Compilation errors** (fixed struct field issues)
7. **Comprehensive validation** - Handler now uses validator for validation
8. **Response data handling** - Tests now get proper response data and headers
9. **Basic validation logic** - 400/404 responses working for most cases

### ❌ Remaining Issues:
1. **Enum validation** - `400_number_enum` test failing (enum validation not working)
2. **Path parameter validation** - `400_invalid_suffix` getting 404 instead of 400
3. **Complex object validation** - Multiple `400_on_dive_*` tests failing (dive validation not working)
4. **Error handling** - `500_Internal_Server_Error` getting 404 instead of 500
5. **JSON unmarshaling errors** - Some tests have JSON parsing issues

## Detailed Work Plan

### Phase 1: Fix Apimodels Structure ✅ COMPLETED
- [x] Add `CreateCookies` struct generation
- [x] Add `Cookies` field to `CreateRequest`
- [x] Fix struct field compilation errors

### Phase 2: Implement Comprehensive Validation ✅ COMPLETED
- [x] Add validator usage in handler method
- [x] Add comprehensive field validation using validator
- [x] Fix validation error responses for basic cases

### Phase 3: Fix Response Data Handling ✅ COMPLETED
- [x] Fix response data structure in generated code
- [x] Implement proper response writing logic
- [x] Add response headers handling

### Phase 4: Fix Advanced Validation Issues (Priority: HIGH)
**Problem**: Advanced validation features not working:
- Enum validation not working
- Path parameter validation (suffix) not implemented
- Complex object validation (dive) not working

**Solution**:
1. **Fix enum validation** - Add proper enum validation tags and logic
2. **Implement path parameter validation** - Add suffix parameter validation
3. **Fix complex object validation** - Implement dive validation for nested objects
4. **Fix error handling** - Ensure proper status codes for different error cases

### Phase 4: Implement Content-Type Checking (Priority: LOW)
**Problem**: Missing content-type validation

**Solution**:
1. **Add content-type switch** statement
2. **Implement proper error handling** for unsupported content types

### Phase 5: Testing and Validation (Priority: HIGH)
**Solution**:
1. **Run tests** after each phase
2. **Fix remaining issues** iteratively
3. **Verify all test cases pass**

## Current Implementation Status

The AST builder now generates:
- ✅ Proper handler structure with validator and create handler
- ✅ Basic parameter parsing (cookies, headers, request body)
- ✅ Proper apimodels structure with RequestCookies
- ✅ Fixed compilation errors

**Next**: Implement comprehensive validation logic using the validator instance.

## Test Data Expectations

The test data expects:
- Complex validation using `validator.Validate` instance
- Proper error responses for validation failures (400/404)
- Response data structure matching test expectations
- Path parameter validation (suffix parameter)
- Cookie validation with proper error handling
- Header validation with time parsing
- Request body validation with comprehensive field checks
