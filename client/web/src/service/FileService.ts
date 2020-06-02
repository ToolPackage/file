
import Axios from './Axios';
import API from './API';

export interface FileInfo {
    fileId: string
    fileName: string
    filePath: string
    fileSize: number
}

export async function getFileList(): Promise<FileInfo[]> {
    let rep = await Axios.get(API.getFileList)

    if (rep.status === 200) {
        return rep.data as FileInfo[]
    } else {
        console.error(rep)
    }
}

export default {
    getFileList
}