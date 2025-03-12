import { z } from 'zod'

export const certificationSchema = z.object({
  tokenId: z.number(),
  owner: z.string(),
  studentId: z.number(),
  certificationType: z.enum(['Diploma', 'Performance']),
  tokenURI: z.string().url(),
  parentId: z.number().optional(),
})

export const certificationListSchema = z.array(certificationSchema)

export type Certification = z.infer<typeof certificationSchema>
