
import Axios from './Axios';
import API from './API';

export interface FileInfo {
    fileId: string
    fileName: string
    filePath: string
    fileSize: number
}

export async function searchFiles() {

}

/**
 * get all file metadata
 */
export async function getFileList(): Promise<FileInfo[]> {
    let rep = await Axios.get(API.getFileList)

    if (rep.status === 200) {
        return rep.data as FileInfo[]
    } else {
        console.error(rep)
    }
}

/**
 * delete file
 * @param fileIds
 */
export async function deleteFile(fileIds: string[]) {

}

interface AsyncTask {

}

class AsyncTaskManager {
    private tasks: Map<string, AsyncTask>

    constructor() {
        this.tasks = new Map()
        // open websocket connection
    }
}

let asyncTaskMgr: AsyncTaskManager

/**
 * post new file
 */
export async function uploadFile() {

}

/**
 * download file
 * @param fileId
 */
export async function downloadFile(fileId: string) {
    
}

export default {
    searchFiles,
    getFileList,
    uploadFile,
    downloadFile,
    deleteFile,
}