import winston, {format} from 'winston'
import {LOG_LEVEL} from '../config'

export function throwIfNullOrUndefined<T = any>(
  value: T | null | undefined,
  msg?: string,
): value is T {
  if (value === null || value === undefined) {
    throw new Error(msg || 'value is not defined')
  }
  return true
}

// There are numerous libraries in npm that can sleep,
// however installing a few of these didn't work in our docker CI environment
// so we will just define the function here.
export function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms)
  })
}

export const logger = winston.createLogger({
  transports: [
    new winston.transports.Console({
      level: LOG_LEVEL,
      format: format.prettyPrint({colorize: true}),
    }),
  ],
})
