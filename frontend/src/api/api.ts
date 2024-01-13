import type { Address } from './entities/address';
import axios from 'axios'
import { LatLng } from 'leaflet'

const api = axios.create({
  baseURL: import.meta.env.API_URL ?? 'http://localhost:3000'
})

export async function fetchAddresses(query: string): Promise<Address[]> {
  const res = await api.get(`/api/search?q=${query}`)
  return res.data as Address[]
}

export async function fetchLocateAddress(address: Address): Promise<LatLng> {
  const res = await api.get(`/api/locate?id=${address.OsmID}`)
  return new LatLng(res.data[1], res.data[0])
}

export async function fetchRoute(points: LatLng[]): Promise<{ geojson:any, time:number, distance:number }[]> {
  const pointString = encodeURIComponent(btoa(JSON.stringify(points.map((p) => [p.lng, p.lat]))));
  const res = await api.get(`/api/route?r=${pointString}`)
  return res.data
}