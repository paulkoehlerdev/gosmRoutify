import type { Address } from './entities/address';
import axios from 'axios'
import type { Point } from '@/api/entities/point'

const api = axios.create({
  baseURL: 'http://localhost:3000/'
})

export async function fetchAddresses(query: string): Promise<Address[]> {
  const res = await api.get(`/api/search?q=${query}`)
  return res.data as Address[]
}

export async function fetchLocateAddress(address: Address): Promise<Point> {
  const res = await api.get(`/api/locate?id=${address.OsmID}`)
  return res.data as Point
}