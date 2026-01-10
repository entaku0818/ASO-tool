'use client'

import { useState } from 'react'
import { setASCCredentials, validateASCCredentials, deleteASCCredentials, ASCCredentials } from '@/lib/api'

type Props = {
  appId: string
  existingCredentials?: ASCCredentials | null
  onSuccess?: () => void
}

export function ASCCredentialsForm({ appId, existingCredentials, onSuccess }: Props) {
  const [issuerId, setIssuerId] = useState(existingCredentials?.issuer_id || '')
  const [keyId, setKeyId] = useState(existingCredentials?.key_id || '')
  const [privateKeyFile, setPrivateKeyFile] = useState<File | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [validationStatus, setValidationStatus] = useState<'idle' | 'validating' | 'valid' | 'invalid'>('idle')
  const [error, setError] = useState<string | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setPrivateKeyFile(file)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!issuerId || !keyId) {
      setError('Issuer ID and Key ID are required')
      return
    }

    if (!existingCredentials && !privateKeyFile) {
      setError('Private key file is required')
      return
    }

    setIsSubmitting(true)
    try {
      let privateKeyBase64 = ''
      if (privateKeyFile) {
        const arrayBuffer = await privateKeyFile.arrayBuffer()
        const bytes = new Uint8Array(arrayBuffer)
        let binary = ''
        bytes.forEach(byte => binary += String.fromCharCode(byte))
        privateKeyBase64 = btoa(binary)
      }

      await setASCCredentials(appId, {
        issuer_id: issuerId,
        key_id: keyId,
        private_key: privateKeyBase64 || 'keep_existing', // Backend will ignore if empty
      })

      onSuccess?.()
    } catch {
      setError('Failed to save credentials')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleValidate = async () => {
    setValidationStatus('validating')
    try {
      const result = await validateASCCredentials(appId)
      setValidationStatus(result.valid ? 'valid' : 'invalid')
    } catch {
      setValidationStatus('invalid')
    }
  }

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete the credentials?')) {
      return
    }

    setIsSubmitting(true)
    try {
      await deleteASCCredentials(appId)
      onSuccess?.()
    } catch {
      setError('Failed to delete credentials')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Issuer ID
        </label>
        <input
          type="text"
          value={issuerId}
          onChange={(e) => setIssuerId(e.target.value)}
          placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
          className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Key ID
        </label>
        <input
          type="text"
          value={keyId}
          onChange={(e) => setKeyId(e.target.value)}
          placeholder="XXXXXXXXXX"
          className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Private Key (.p8 file)
        </label>
        <input
          type="file"
          accept=".p8"
          onChange={handleFileChange}
          className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {existingCredentials && !privateKeyFile && (
          <p className="text-xs text-gray-500 mt-1">Leave empty to keep existing key</p>
        )}
      </div>

      {error && (
        <div className="text-red-600 text-sm">{error}</div>
      )}

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={isSubmitting}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
        >
          {isSubmitting ? 'Saving...' : 'Save'}
        </button>

        {existingCredentials && (
          <>
            <button
              type="button"
              onClick={handleValidate}
              disabled={validationStatus === 'validating'}
              className="px-4 py-2 border rounded hover:bg-gray-50"
            >
              {validationStatus === 'validating' ? 'Validating...' : 'Validate'}
            </button>
            <button
              type="button"
              onClick={handleDelete}
              disabled={isSubmitting}
              className="px-4 py-2 text-red-600 border border-red-300 rounded hover:bg-red-50"
            >
              Delete
            </button>
          </>
        )}
      </div>

      {validationStatus === 'valid' && (
        <div className="text-green-600 text-sm">Credentials are valid</div>
      )}
      {validationStatus === 'invalid' && (
        <div className="text-red-600 text-sm">Credentials are invalid</div>
      )}
    </form>
  )
}
