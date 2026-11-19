import assert from 'node:assert/strict'
import { test } from 'node:test'
import { safeRedirect } from './redirect.js'

test('keeps the invitation query and fragment when returning to the website', () => {
  assert.equal(safeRedirect('/invite?aff=A%2BB#share', 'https://site.example'), '/invite?aff=A%2BB#share')
})

test('rejects external URLs and authentication redirect loops', () => {
  for (const value of [undefined, ['bad'], '//other.example', '/\\other.example', 'https://other.example', '/admin', '/admin/login', '/login?redirect=/', '/register']) {
    assert.equal(safeRedirect(value, 'https://site.example'), '/')
  }
})
