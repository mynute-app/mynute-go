-- Fix GetEmployeeWorkRange controller name mismatch
UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id'
WHERE controller_name = 'GetEmployeeWorkRangeById'
  AND method = 'GET'
  AND path = '/employee/:id/work_range/:work_range_id';
