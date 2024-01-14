import type { Address } from './entities/address';
import { LatLng } from 'leaflet'

export async function fetchAddresses(baseURL: string, query: string): Promise<Address[]> {
  const res = await fetch(`${baseURL}/api/search?q=${query}`)
  return await res.json() as Address[]
}

export async function fetchLocateAddress(baseURL: string, address: Address): Promise<LatLng> {
  const res = await fetch(`${baseURL}/api/locate?id=${address.OsmID}`)
  const json = await res.json()
  return new LatLng(json[1], json[0])
}

export async function fetchRoute(baseURL: string, points: LatLng[]): Promise<{ geojson:any, time:number, distance:number }[]> {
  const pointString = encodeURIComponent(btoa(JSON.stringify(points.map((p) => [p.lng, p.lat]))));
  const res = await fetch(`${baseURL}/api/route?r=${pointString}`)
  const json = await res.json()
  return json
}