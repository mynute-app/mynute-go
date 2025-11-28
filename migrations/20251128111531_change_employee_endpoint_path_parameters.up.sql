-- Migration: change_employee_endpoint_path_parameters
-- Created at: 20251128111531
-- Description: Updates endpoint paths from generic :id to specific :employee_id and :client_id
--              to match the current code definitions and ensure authorization works correctly
--              This fixes the parameter naming standardization done in commit 4afebdf

-- Employee Endpoints: Update :id to :employee_id where applicable
UPDATE public.endpoints
SET path = '/employee/:employee_id/work_schedule'
WHERE method = 'POST' 
  AND path = '/employee/:id/work_schedule'
  AND controller_name = 'CreateEmployeeWorkSchedule';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_schedule'
WHERE method = 'GET' 
  AND path = '/employee/:id/work_schedule'
  AND controller_name = 'GetEmployeeWorkSchedule';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id'
WHERE method = 'GET' 
  AND path = '/employee/:id/work_range/:work_range_id'
  AND controller_name = 'GetEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id'
WHERE method = 'DELETE' 
  AND path = '/employee/:id/work_range/:work_range_id'
  AND controller_name = 'DeleteEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id'
WHERE method = 'PUT' 
  AND path = '/employee/:id/work_range/:work_range_id'
  AND controller_name = 'UpdateEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id/services'
WHERE method = 'POST' 
  AND path = '/employee/:id/work_range/:work_range_id/services'
  AND controller_name = 'AddEmployeeWorkRangeServices';

UPDATE public.endpoints
SET path = '/employee/:employee_id/work_range/:work_range_id/service/:service_id'
WHERE method = 'DELETE' 
  AND path = '/employee/:id/work_range/:work_range_id/service/:service_id'
  AND controller_name = 'DeleteEmployeeWorkRangeService';

UPDATE public.endpoints
SET path = '/employee/:employee_id/appointments'
WHERE method = 'GET' 
  AND path = '/employee/:id/appointments'
  AND controller_name = 'GetEmployeeAppointmentsById';

-- Note: GetEmployeeById, UpdateEmployeeById, DeleteEmployeeById, UpdateEmployeeImages, 
-- and DeleteEmployeeImage still use :id in the code, so we don't update those
