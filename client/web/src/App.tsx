
import React, { Component } from 'react';
import { Container, Divider, Button, Input, Icon, Form, List, Progress, Header } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css?global';
import { initializeFileTypeIcons, getFileTypeIconProps } from '@uifabric/file-type-icons';
import FileService, { FileInfo } from './service/FileService';
import {Icon as FluentuiIcon} from '@fluentui/react/lib/Icon';
import './App.css';

initializeFileTypeIcons()

interface AppState {
    fileList: FileInfo[]
    operationList: JSX.Element[]
}

export default class App extends Component<any, AppState> {

    constructor(props: any) {
        super(props)
        this.state = {
            fileList: [{fileName: 'test_file.docx', fileId: 'asd', fileSize: 3, filePath: ''}],
            operationList: [
                <OperationProgress fileName='test file' fileId='xxx' type='download' />,
                <OperationProgress fileName='test file' fileId='xxx' type='upload' />
            ]
        }
    }

    componentDidMount() {
        FileService.getFileList().then(fileList => this.setState({fileList}))
    }

    render() {
        const { fileList, operationList } = this.state

        return (
            <div className='App'>
                <Container style={{position: 'relative', width: '40%', paddingTop: 50}}>
                    {/* header */}
                    <Form>
                        <Form.Group inline>
                            <Form.Field width='13'>
                                <Input size='mini' action fluid placeholder='Search...'>
                                    <input />
                                    <Button size='mini' icon='search' primary />
                                </Input>
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='upload' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='download' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='setting' primary />
                            </Form.Field>
                        </Form.Group>
                    </Form>
                    <Divider />
                    {/* file list */}
                    <List animated>
                        {
                            fileList.map((fileInfo, idx) => (
                                <List.Item key={idx} style={{cursor: 'pointer'}}>
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
                                    <List.Item key={idx}>
                                        {v}
                                    </List.Item>
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
                <div style={{display: 'flex', flexFlow: 'column', width: '100%'}}>
                    <Header size='tiny' style={{display: 'inline-block'}}>{this.props.fileName}</Header>
                    <Progress style={{marginBottom: 0}} percent={33} progress indicating />
                </div>
                <div style={{display: 'flex'}}>
                    <Button icon='cancel'/>
                </div>
            </div>
        )
    }
}