
import React, { Component } from 'react';
import { Container, Divider, Button, Input, Icon, Form, List, Progress, Header, Modal } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css?global';
import { initializeFileTypeIcons, getFileTypeIconProps } from '@uifabric/file-type-icons';
import FileService, { FileInfo } from './service/FileService';
import {Icon as FluentuiIcon} from '@fluentui/react/lib/Icon';
import './App.css';

initializeFileTypeIcons()

interface AppState {
    fileList: FileInfo[]
    selectedItemList: Set<number>
    operationList: JSX.Element[]
    openDeleteWaring: boolean
}

export default class App extends Component<any, AppState> {

    constructor(props: any) {
        super(props)
        this.state = {
            fileList: [{fileName: 'test_file.docx', fileId: 'asd', fileSize: 3, filePath: ''}, {fileName: 'test_file.docx', fileId: 'asd', fileSize: 3, filePath: ''}],
            selectedItemList: new Set(),
            operationList: [
                <OperationProgress fileName='test file' fileId='xxx' type='download' />,
                <OperationProgress fileName='test file' fileId='xxx' type='upload' />
            ],
            openDeleteWaring: false
        }
    }

    componentDidMount() {
        FileService.getFileList().then(fileList => this.setState({fileList}))
    }

    openUploadFileModal() {
        let uploadEntry = this.refs.uploadEntry as HTMLInputElement
        uploadEntry.click()
    }

    selectFile(idx: number) {
        const { selectedItemList } = this.state
        if (selectedItemList.has(idx)) {
            // remove selected
            selectedItemList.delete(idx)
        } else {
            selectedItemList.add(idx)
        }
        this.forceUpdate()
    }

    selectAllFile() {
        const { fileList, selectedItemList } = this.state
        if (selectedItemList.size === fileList.length) {
            // de-select all
            selectedItemList.clear()
        } else {
            fileList.forEach((_, idx) => selectedItemList.add(idx))
        }
        this.forceUpdate()
    }

    getSelectedFileList() {
        const { fileList, selectedItemList } = this.state
        let selectedFileList: FileInfo[] = []
        selectedItemList.forEach((idx) => selectedFileList.push(fileList[idx]))
        return selectedFileList
    }

    deleteSelectedFiles() {
        // TODO:
    }

    render() {
        const { fileList, selectedItemList, operationList,
            openDeleteWaring } = this.state

        return (
            <div className='App'>
                <Container style={{position: 'relative', width: '40%', paddingTop: 50}}>
                    {/* header */}
                    <input ref='uploadEntry' style={{display: 'none'}} type='file' />
                    <Form>
                        <Form.Group inline>
                            <Form.Field width='13'>
                                <Input size='mini' action fluid placeholder='Search...'>
                                    <input />
                                    <Button size='mini' icon='search' primary />
                                </Input>
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='setting' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='check' primary
                                    onClick={this.selectAllFile.bind(this)} />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='download' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='upload' primary
                                    onClick={this.openUploadFileModal.bind(this)} />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='delete' color='red'
                                    onClick={() => {
                                        if (selectedItemList.size > 0) {
                                            this.setState({openDeleteWaring: true})
                                        } else {
                                            // TODO: warning no file selected
                                        }
                                    }}/>
                            </Form.Field>
                        </Form.Group>
                    </Form>
                    {/* delete warning */}
                    <Modal open={openDeleteWaring}>
                        <Modal.Header style={{color: 'red'}}>Operation Warning</Modal.Header>
                        <Modal.Content image>
                            <Modal.Description>
                                <Header>Are you sure to delete the following files?</Header>
                                <List>
                                    { this.getSelectedFileList().map((file, idx) => (
                                        <List.Item key={idx}>
                                            {file.fileName}
                                        </List.Item>
                                    ))}
                                </List>
                            </Modal.Description>
                        </Modal.Content>
                        <Modal.Actions>
                            <Button color='red' onClick={() => this.setState({openDeleteWaring: false})}>Cancel</Button>
                            <Button primary onClick={this.deleteSelectedFiles.bind(this)}>Confirm</Button>
                        </Modal.Actions>
                    </Modal>
                    <Divider />
                    {/* file list */}
                    <List>
                        {
                            fileList.map((fileInfo, idx) => (
                                <List.Item className={ 'file' + (selectedItemList.has(idx) ? ' selectedFile' : '') } key={idx} style={{cursor: 'pointer'}}
                                    onClick={() => this.selectFile(idx)}>
                                    <List.Icon verticalAlign='middle'>
                                        <FluentuiIcon {...getFileTypeIconProps({extension: 'docx', size: 16})} />
                                    </List.Icon>
                                    <List.Content>
                                        <List.Header>{fileInfo.fileName}</List.Header>
                                        <List.Description>
                                            <span>{fileInfo.fileSize}</span>
                                            <span style={{marginLeft: 10}}>Uploaded 3 mins ago</span>
                                        </List.Description>
                                    </List.Content>
                                </List.Item>))
                        }
                    </List>
                    {/* upload & download progress */}
                    <div className='progressContainer'>
                        <List className=''>
                            {
                                operationList.map((v, idx) => (
                                    <List.Item key={idx}>{v}</List.Item>
                                ))
                            }
                        </List>
                    </div>
                </Container>
            </div>
        )
    }
}

interface OperationProgressProps {
    fileName: string
    fileId: string
    type: 'upload' | 'download'
}

class OperationProgress extends Component<OperationProgressProps> {

    render() {
        return (
            <div className='opeartionProgress'>
                <div style={{display: 'flex', flexFlow: 'column', justifyContent: 'flex-end'}}>
                    <Icon name={this.props.type} style={{color: '#2185d0'}} size='big' />
                </div>
                <div style={{display: 'flex', flexFlow: 'column', width: '100%', marginRight: '4px'}}>
                    <Header size='tiny' style={{display: 'inline-block'}}>{this.props.fileName}</Header>
                    <Progress style={{marginBottom: 0}} percent={33} progress indicating />
                </div>
                <div style={{display: 'flex', flexFlow: 'column', justifyContent: 'flex-end'}}>
                    <Button icon='cancel' size='mini' color='red' />
                </div>
            </div>
        )
    }
}