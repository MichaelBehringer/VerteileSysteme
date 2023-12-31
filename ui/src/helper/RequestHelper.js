import axios from "axios";

//Load enpdoint url dynamic from url of the browser
export const baseUrl = window.location.hostname
const url = "http://"+baseUrl+"/lobby/"

export async function doPostRequest(path, param) {
	return axios.post(url+path, param)
}

export async function doCustomPostRequest(path, param) {
	return axios.post("http://"+baseUrl+"/"+path, param)
}

export async function doPostRequestAuth(path, param, auth) {
	return axios.post(url+path, param, {headers: {Authorization: 'Bearer ' + auth}})
}

export async function doGetRequest(path) {
	return axios.get(url+path)
}

export async function doGetRequestParam(path, param) {
	return axios.get(url+path+"/"+param)
}

export async function doGetRequestBlob(path) {
	return axios.get(url+path, { responseType: 'blob' })
}

export async function doGetRequestAuth(path, auth) {
	return axios.get(url+path, {headers: {Authorization: 'Bearer ' + auth}})
}

export async function doDeleteRequest(path, param) {
	return axios.delete(url+path, param)
}

export async function doDeleteRequestAuth(path, param, auth) {
	const dataObj = { data: param, headers: {Authorization: 'Bearer ' + auth}}
	return axios.delete(url+path, dataObj)
}

export async function doPutRequest(path, param) {
	return axios.put(url+path, param)
}

export async function doPutRequestAuth(path, param, auth) {
	return axios.put(url+path, param, {headers: {Authorization: 'Bearer ' + auth}})
}