export default ({
    getFileList: '/files',
    uploadFile: '/files',
    downloadFile: (fileId: number) => `/files/${fileId}`,
    deleteFile: (fileId: number) => `/files/${fileId}`
})