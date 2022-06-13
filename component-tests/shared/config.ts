export const SERVICE_HOST =
  process.env['COMPONENT_TESTS_PERSONAL_DATA_HOST'] || 'http://localhost:8057'
export const MONGO_DSN: string = process.env['MONGO_DSN'] || 'TBD'
export const COMPONENT_TESTS_WIREMOCK_URL =
  process.env['COMPONENT_TESTS_WIREMOCK_URL'] || 'http://localhost:8581'
export const COMPONENT_TEST_TIMEOUT = process.env['COMPONENT_TEST_TIMEOUT'] || 15
export const LOG_LEVEL = 'info'
