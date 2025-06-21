import { HttpException, HttpStatus, Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { CozeEntity } from './coze.entity';
import { CozeChatDto, CozeCreateDto, CozeListDto, CozeUpdateDto } from './dto/coze.dto';
import { CozeAPI, COZE_CN_BASE_URL, RoleType, ChatStatus, getJWTToken, ChatEventType } from '@coze/api';
import { HttpService } from '@nestjs/axios';
import axios from 'axios';

@Injectable()
export class CozeService {
  private logger = new Logger(CozeService.name);
  private client: CozeAPI;
  private accessToken: string | null = null;
  
  // 默认配置
  private defaultConfig = {
    coze_api_base: "https://api.coze.cn",
    coze_client_id: "1123962302922",
    coze_public_key_id: "fxkr4uRho76_yAqNKA_1wV7on2AwjQIk3tjkzsnk-Z4",
    coze_private_key: "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCttH3mlTJKKhQT\ny8RJBVIZDXKBdJSx9Y2VSFkmEK6ASMvE5W4VW7qX41PMg5VtGlajQCWxhugJwV1H\npEKvKqU0RsIZ14CTgTuVekRiruCnf0DN6KM6WuEJNsQxCTICchFC4sOt7kyjtCHU\nYUlF4faAFqhyMLzjJpdss+QQngSEjH1s+alr+RySBSA5puEF0VFMpRMEjFwi4Jno\n0yPNFCOiYvz+tRIK6OhSdtfhMKAmmt+62/Op3v9O2GAhFVJO1JNl0Odxlf/rs3nB\nv+37ds6pz+7LDOzGSufAAoigMDRRn1a93AT2Q4tlz9kYe5RncKhnA6xPgomNZfIf\nVMwjOGF1AgMBAAECggEAARbfr0GCRjrLU3B0s6yH3kZaUHuFrzQGBkik3ns+TOmn\n9X0m2pVvryIq1V6B4mRG5NEzK1DYRa9jwV5DWMvgq1pCP109ni8yS3av1RqZqBNB\nOclatLP7M06XnmMbYC6M8ylu5rlW27P2fll51ylanWUG+2hY1ufYDUN3i68iAh7I\nwsMIPdf1+TYGi8IdHg2hHpBEMecMU/bhRaDXM3DOIMWYZ8l6k7uxrsDKv7Uhov+n\nZIfv+AjWqz7N+yvIR6GM0HqG+iv4B+JAkrcYPbchl44IPbbo4HKRjEY2+LkOl/5I\nwM46pO6eAuUNXhXdw0y3WRqOdc/o+c7PD/JeZWTgpQKBgQDyu7hsKWapbpeZ5Kqm\nBxr+yec0/xcfXxkhcc9Ehdd8+v/QGKza8vRkCUp9elpccMpMW5rLdPVZpP2UovXF\nrSjDKl9m2H2OKzzBRO8uT/a//XSzywlQjcmSaYmu0stBhf7tSoDjlUjFqvYjdxDO\n1u2cqZiymVyC3R/bQv1RPm08zwKBgQC3MvDjBWeV6YBcu/iSF8jmokSOz71YyUFn\n9CReKzQzEDO6DVhWulsIDlPne3dfqI2vST+ZljNUo1S8jxaMX8C2rTdKmIGVWOQP\ntbQftdtti5QJ83Z5HhE4KSZ6irq1GY/AmL1DiU2mRHch7KOMKDGrTps1Rwv/Y8t2\n72Ha9q+2ewKBgQC88JwELVHJDtmYo6Kla6B6tSRwXyNbewWvv8wLVXc/xIy9KYfb\nQgQznfvKoiOWEwGU4DUkq5yTM9djDFnsjfXNvLzX7CoHMOawtfzLetjh5uMhVCii\n+Erv2ZCfcVtfXHLrt/ONstUbcBD52CNQLYJ1UJoYY0HcZ0z1ujY+OC6FhwKBgCpx\n0v3GMsm439SceGrgt9s3nUq5NtVrS4waNJLcz6tFBbcFgIIXix/Csg3fvTichLcn\n8WRUOHBTpz5IqKC9TpkEaNsPmnZPsgcxwhnWuJAY1qO3lKtbHAI3BoM9wSRUV8n3\nmWIcXbE4C6IAgaPnbBqUi8E8RLtXE7zqmXFx1iQhAoGAS0OlIv/Nr7K0cWdQouAT\nO7ohkipKJBKgD6rdU3e2lOtl3vheW7hxipe7id5hMwibeplbxEco1voxrn2qHiIv\n+g/WoSEOUkukzG7hyrlHL26L5Lpr3pTiTLpw3G3edyZqgR+mOIweEmQtQxIyOuyP\nepWwT2VHi5YaRSYwKgklrPk=\n-----END PRIVATE KEY-----",
    default_bot_id: "7483341266928468008",
    space_id: "7477025304549294118"
  };

  constructor(
    @InjectRepository(CozeEntity)
    private readonly cozeRepository: Repository<CozeEntity>,
    private readonly httpService: HttpService
  ) {
    this.logger.log('CozeService初始化中...');
    try {
      this.initCozeClient();
      this.logger.log('CozeService初始化完成');
    } catch (error) {
      this.logger.error(`CozeService初始化失败: ${error.message}`, error.stack);
      // 初始化失败但不中断服务启动，而是记录错误
    }
  }

  // 初始化Coze客户端
  private initCozeClient() {
    try {
      this.logger.log('开始初始化CozeAPI客户端');
      
      // 创建获取token的函数
      const getToken = async () => {
        try {
          this.logger.log('尝试获取JWT Token');
          const token = await getJWTToken({
            baseURL: COZE_CN_BASE_URL,
            appId: this.defaultConfig.coze_client_id,
            aud: 'api.coze.cn',
            keyid: this.defaultConfig.coze_public_key_id,
            privateKey: this.defaultConfig.coze_private_key,
          });
          
          if (token && token.access_token) {
            this.logger.log(`JWT Token获取成功: ${token.access_token.substring(0, 10)}...`);
            this.accessToken = token.access_token;
            return token.access_token;
          } else {
            this.logger.error('JWT Token获取失败: 返回结果为空或无效');
            throw new Error('Invalid token response');
          }
        } catch (tokenError) {
          this.logger.error(`JWT Token获取失败: ${tokenError.message}`, tokenError.stack);
          throw tokenError;
        }
      };

      // 验证一下token是否能正常获取
      this.logger.log('验证Token获取功能');
      getToken().then(token => {
        this.logger.log('Token验证成功，继续初始化');
      }).catch(error => {
        this.logger.error(`Token验证失败，但将继续尝试初始化: ${error.message}`);
      });

      this.logger.log('初始化CozeAPI实例...');
      
      // 初始化CozeAPI客户端
      this.client = new CozeAPI({
        token: getToken,
        baseURL: COZE_CN_BASE_URL,
      });

      // 检查client是否正确初始化
      if (this.client) {
        this.logger.log('CozeAPI客户端创建成功');
        // 检查client的属性和方法
        const clientAny = this.client as any;
        this.logger.log(`客户端信息: workflows=${!!clientAny.workflows}, runs=${!!clientAny.runs}`);
        if (clientAny.workflows) {
          this.logger.log(`可用的workflows方法: ${Object.keys(clientAny.workflows).join(', ')}`);
        }
      } else {
        this.logger.warn('CozeAPI客户端创建可能不完整');
      }

      this.logger.log('CozeAPI客户端初始化成功');
    } catch (error) {
      this.logger.error(`CozeAPI客户端初始化失败: ${error.message}`);
      this.logger.error(`初始化失败详情: ${error.stack}`);
      // 重新抛出错误以便上层catch捕获
      throw error;
    }
  }

  // 创建Coze配置
  async create(cozeCreateDto: CozeCreateDto) {
    try {
      // 检查是否已存在相同的workflowId
      const exists = await this.cozeRepository.findOne({
        where: { workflowId: cozeCreateDto.workflowId },
      });

      if (exists) {
        throw new HttpException('Workflow ID已存在', HttpStatus.BAD_REQUEST);
      }

      // 创建新的Coze配置
      const coze = this.cozeRepository.create(cozeCreateDto);
      return await this.cozeRepository.save(coze);
    } catch (error) {
      this.logger.error(`创建Coze配置失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '创建Coze配置失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 更新Coze配置
  async update(id: number, cozeUpdateDto: CozeUpdateDto) {
    try {
      // 查找要更新的记录
      const coze = await this.cozeRepository.findOne({ where: { id } });
      if (!coze) {
        throw new HttpException('Coze配置不存在', HttpStatus.NOT_FOUND);
      }

      // 如果更新 workflowId，检查是否与其他记录重复
      if (cozeUpdateDto.workflowId && cozeUpdateDto.workflowId !== coze.workflowId) {
        const exists = await this.cozeRepository.findOne({
          where: { workflowId: cozeUpdateDto.workflowId },
        });
        
        if (exists && exists.id !== id) {
          throw new HttpException('Workflow ID已被其他配置使用', HttpStatus.BAD_REQUEST);
        }
      }

      // 合并并保存
      const updatedCoze = Object.assign(coze, cozeUpdateDto);
      return await this.cozeRepository.save(updatedCoze);
    } catch (error) {
      this.logger.error(`更新Coze配置失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '更新Coze配置失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 删除Coze配置
  async delete(id: number) {
    try {
      const coze = await this.cozeRepository.findOne({ where: { id } });
      if (!coze) {
        throw new HttpException('Coze配置不存在', HttpStatus.NOT_FOUND);
      }
      
      await this.cozeRepository.softDelete(id);
      return { success: true };
    } catch (error) {
      this.logger.error(`删除Coze配置失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '删除Coze配置失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 获取Coze配置列表
  async findAll(query: CozeListDto) {
    try {
      this.logger.log(`Getting Coze list with query: ${JSON.stringify(query)}`);
      
      const page = query.page || 1;
      const size = query.size || 10;
      
      const [items, total] = await this.cozeRepository.findAndCount({
        skip: (page - 1) * size,
        take: size,
        order: { createdAt: 'DESC' },
      });

      const result = {
        items,
        total,
        page,
        size,
      };
      
      this.logger.log(`Found ${total} Coze configs`);
      return result;
    } catch (error) {
      this.logger.error(`获取Coze配置列表失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '获取Coze配置列表失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 获取单个Coze配置
  async findOne(id: number) {
    try {
      const coze = await this.cozeRepository.findOne({ where: { id } });
      if (!coze) {
        throw new HttpException('Coze配置不存在', HttpStatus.NOT_FOUND);
      }
      return coze;
    } catch (error) {
      this.logger.error(`获取Coze配置失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '获取Coze配置失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 查找指定WorkflowID的Coze配置
  async findByWorkflowId(workflowId: string) {
    try {
      this.logger.log(`查找workflowId为 ${workflowId} 的Coze配置`);
      
      const coze = await this.cozeRepository.findOne({ 
        where: { workflowId, status: 1 } 
      });
      
      if (!coze) {
        this.logger.warn(`未找到workflowId为 ${workflowId} 的Coze配置`);
        throw new HttpException('Coze Workflow不存在或未启用', HttpStatus.NOT_FOUND);
      }
      
      this.logger.log(`找到workflowId为 ${workflowId} 的Coze配置: ${coze.name}`);
      return coze;
    } catch (error) {
      this.logger.error(`获取Coze Workflow配置失败: ${error.message}`, error.stack);
      throw new HttpException(
        error.message || '获取Coze Workflow配置失败',
        error.status || HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  /**
   * 获取机器人列表
   */
  async getBotList() {
    try {
      this.logger.log('获取机器人列表');
      
      // 确保client已初始化
      if (!this.client) {
        throw new Error('CozeAPI客户端未初始化');
      }
      
      // 使用SDK的原生方法获取机器人列表
      const api = this.client as any;
      if (api.bots && typeof api.bots.list === 'function') {
        const response = await api.bots.list({
          space_id: this.defaultConfig.space_id,
        });
        
        this.logger.log(`获取机器人列表成功, 共${response.bots?.length || 0}个机器人`);
        return response;
      } else {
        throw new Error('SDK中没有可用的机器人列表获取方法');
      }
    } catch (error) {
      this.logger.error(`获取机器人列表失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '获取机器人列表失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  /**
   * 聊天 - 非流式响应
   */
  async chat(chatDto: CozeChatDto) {
    try {
      // 获取Workflow ID
      let workflowId = chatDto.workflowId;
      if (!workflowId) {
        workflowId = this.defaultConfig.default_bot_id; // 暂时使用default_bot_id作为默认值
        this.logger.log(`未指定workflowId，使用默认workflowId: ${workflowId}`);
      }

      this.logger.log(`开始聊天请求, workflowId=${workflowId}, message=${chatDto.message}`);
      
      // 使用SDK的createAndPoll方法进行聊天
      const result = await this.client.chat.createAndPoll({
        bot_id: workflowId, // SDK仍使用bot_id参数
        additional_messages: [{
          role: RoleType.User,
          content: chatDto.message,
          content_type: 'text',
        }],
        user_id: chatDto.sessionId || 'default_user_id',
      });

      this.logger.log(`聊天请求完成, 状态: ${result.chat?.status}`);

      // 检查响应状态
      if (result.chat?.status === ChatStatus.COMPLETED) {
        const assistantMessages = result.messages.filter(msg => 
          msg.role === 'assistant' && msg.content_type === 'text'
        );
        
        if (assistantMessages.length > 0) {
          const latestMessage = assistantMessages[assistantMessages.length - 1];
          return {
            content: latestMessage.content,
            conversation_id: result.chat.id || null,
          };
        }
      }

      throw new HttpException('未获取到有效回复', HttpStatus.BAD_REQUEST);
    } catch (error) {
      this.logger.error(`聊天请求失败: ${error.message}`);
      if (error.response) {
        this.logger.error(`错误响应: ${JSON.stringify(error.response.data || {})}`);
      }
      throw new HttpException(
        error.message || '聊天请求失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
  
  /**
   * 聊天 - 流式响应
   */
  async streamChat(chatDto: CozeChatDto) {
    try {
      // 获取Workflow ID
      let workflowId = chatDto.workflowId;
      if (!workflowId) {
        workflowId = this.defaultConfig.default_bot_id; // 暂时使用default_bot_id作为默认值
        this.logger.log(`未指定workflowId，使用默认workflowId: ${workflowId}`);
      }

      this.logger.log(`开始流式聊天请求, workflowId=${workflowId}, message=${chatDto.message}`);
      
      return this.client.chat.stream({
        bot_id: workflowId, // SDK仍使用bot_id参数
        additional_messages: [{
          role: RoleType.User,
          content: chatDto.message,
          content_type: 'text',
        }],
        user_id: chatDto.sessionId || 'default_user_id',
      });
    } catch (error) {
      this.logger.error(`流式聊天请求失败: ${error.message}`);
      throw new HttpException(
        error.message || '流式聊天请求失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  /**
   * 获取工作流列表 - 使用@coze/api
   */
  async getWorkflowList(spaceId?: string) {
    try {
      this.logger.log('获取工作流列表');
      
      // 确保client已初始化
      if (!this.client || !(this.client as any).workflows) {
        throw new Error('CozeAPI客户端未初始化或workflows不可用');
      }
      
      // 使用CozeAPI client
      const targetSpaceId = spaceId || this.defaultConfig.space_id;
      this.logger.log(`工作流列表使用的spaceId: ${targetSpaceId}`);
      
      // 使用@coze/api的方法
      const api = this.client as any;
      
      // 记录当前client实例的信息
      this.logger.log(`客户端信息: ${JSON.stringify({
        hasWorkflows: !!api.workflows,
        workflowsType: api.workflows ? typeof api.workflows : 'undefined',
        methods: api.workflows ? Object.keys(api.workflows).join(',') : 'none'
      })}`);
      
      // 尝试刷新token
      try {
        // 使用dynamicAny来避免类型错误
        this.logger.log('尝试刷新token');
        const dynamicClient = this.client as any;
        if (typeof dynamicClient.token === 'function') {
          await dynamicClient.token();
          this.logger.log('成功获取token');
        } else {
          this.logger.log('token方法不可用');
        }
      } catch (tokenError) {
        this.logger.warn(`获取token失败: ${tokenError.message}`);
      }
      
      // 尝试listPublished方法
      this.logger.log('尝试使用listPublished方法');
      try {
        const response = await api.workflows.listPublished({
          space_id: targetSpaceId
        });
        
        this.logger.log(`获取工作流列表成功, 共${response.workflows?.length || 0}个工作流`);
        return response;
      } catch (methodError) {
        this.logger.warn(`listPublished方法调用失败: ${methodError.message}`);
        
        // 尝试替代方法
        this.logger.log('尝试使用替代方法');
        if (api.workflows.published && typeof api.workflows.published.list === 'function') {
          this.logger.log('发现published.list方法，尝试调用');
          const response = await api.workflows.published.list({
            space_id: targetSpaceId
          });
          
          this.logger.log(`获取工作流列表成功, 共${response.workflows?.length || 0}个工作流`);
          return response;
        } else {
          throw new Error('未找到可用的SDK工作流列表获取方法');
        }
      }
    } catch (error) {
      this.logger.error(`获取工作流列表失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '获取工作流列表失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  /**
   * 处理工作流参数，特别是文件类型的参数
   * @param parameters 原始参数
   * @returns 处理后符合Coze API格式的参数
   */
  private processWorkflowParameters(parameters?: Record<string, any>): Record<string, any> {
    const processedParameters = {};
    
    if (!parameters) {
      return processedParameters;
    }
    
    for (const [key, value] of Object.entries(parameters)) {
      if (key.startsWith('file_') || key.includes('_file_')) {
        // 处理文件类型参数
        if (Array.isArray(value)) {
          this.logger.log(`处理多文件参数元数据 ${key}: ${Array.isArray(value)}`);
          // 多文件情况: 转换为[{"file_id": "xxx"}, {"file_id": "yyy"}]格式的字符串
          const fileObjects = value.map(fileId => ({ file_id: fileId }));
          processedParameters[key] = JSON.stringify(fileObjects);
          this.logger.log(`处理多文件参数 ${key}: ${processedParameters[key]}`);
        } else {
          // 单文件情况: 转换为{"file_id": "xxx"}格式的字符串
          processedParameters[key] = JSON.stringify({ file_id: value });
          this.logger.log(`处理单文件参数 ${key}: ${processedParameters[key]}`);
        }
      } else {
        // 非文件类型参数直接传递
        processedParameters[key] = value;
      }
    }
    
    this.logger.log(`工作流参数处理完成: ${JSON.stringify(processedParameters)}`);
    return processedParameters;
  }

  /**
   * 使用@coze/api触发工作流
   * 用户可以在parameters中直接传入is_async和callback_url来控制执行方式
   * @param parameters 工作流参数，可包含is_async和callback_url
   * @returns 工作流执行结果
   */
  async runWorkflow(parameters: Record<string, any>) {
    try {
      // 提取工作流ID
      const workflowId = parameters.workflow_id;
      
      // 验证必要参数
      if (!workflowId) {
        throw new HttpException(
          '缺少必要的工作流ID', 
          HttpStatus.BAD_REQUEST
        );
      }
      
      this.logger.log(
        `触发工作流, workflowId=${workflowId}` + 
        (parameters.user_id ? `, userId=${parameters.user_id}` : '')
      );
      
      // 处理参数 - 过滤控制参数
      const processedParams = this.processWorkflowParameters(parameters);
      
      // 使用SDK的原生方法触发工作流
      const api = this.client as any;
      const response = await api.workflows.runs.create({
        workflow_id: workflowId,
        is_async: parameters.is_async,
        parameters: processedParams,
        user_id: parameters.user_id || 'default_user_id',
      });

      this.logger.log(`工作流触发成功, is_async=${parameters.is_async}`);
      this.logger.log(`工作流触发成功, response=${JSON.stringify(response)}`);
      this.logger.log(`工作流触发成功, executeId=${response.execute_id}`);
      return response;
    } catch (error) {
      this.logger.error(`触发工作流失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '触发工作流失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
  
  /**
   * 流式执行工作流 - 同样支持参数中指定workflow_id
   */
  async streamWorkflow(parameters: Record<string, any>) {
    try {
      // 处理参数
      const processedParams = this.processWorkflowParameters(parameters);
      
      // 获取工作流ID
      const workflowId = parameters.workflow_id;
      
      // 验证必要参数
      if (!workflowId) {
        throw new HttpException('缺少必要的工作流ID', HttpStatus.BAD_REQUEST);
      }
      
      this.logger.log(`流式触发工作流, workflowId=${workflowId}`);
      
      // 使用SDK触发工作流
      const api = this.client as any;
      const result = await api.workflows.runs.stream({
        workflow_id: workflowId,
        parameters: processedParams,
        user_id: parameters.user_id || 'default_user_id',
      });
      
      return result;
    } catch (error) {
      this.logger.error(`流式触发工作流失败: ${error.message}`);
      throw new HttpException(
        error.message || '流式触发工作流失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
  
  /**
   * 获取工作流执行历史 - 使用@coze/api
   */
  async getWorkflowHistory(workflowId: string, executeId: string) {
    try {
      this.logger.log(`获取工作流历史, workflowId=${workflowId}, executeId=${executeId}`);
      
      // 使用@coze/api的原生方法获取工作流执行历史
      const history = await this.client.workflows.runs.history(
        workflowId,
        executeId
      );
      
      this.logger.log(`获取工作流历史成功`);
      return history;
    } catch (error) {
      this.logger.error(`获取工作流历史失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '获取工作流历史失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  
  
  /**
   * 按状态获取工作流执行历史 - 使用@coze/api
   * 参考文档: https://www.coze.cn/open/docs/developer_guides/workflow_history
   */
  async getWorkflowHistoryByStatus(workflowId: string, status?: string, page: number = 1, pageSize: number = 10) {
    try {
      this.logger.log(`获取工作流执行历史, workflowId=${workflowId}, status=${status || '全部'}, page=${page}, pageSize=${pageSize}`);
      
      // 确保client已初始化
      if (!this.client || !(this.client as any).workflows) {
        throw new Error('CozeAPI客户端未初始化或workflows不可用');
      }
      
      // 准备请求参数
      const requestParams: any = {
        workflow_id: workflowId,
        page,
        page_size: pageSize
      };
      
      if (status) {
        requestParams.status = status;
      }
      
      // 使用@coze/api的原生方法获取执行历史
      const api = this.client as any;
      let result;
      
      if (api.workflows.runs && typeof api.workflows.runs.listHistory === 'function') {
        this.logger.log('使用runs.listHistory方法获取工作流执行历史');
        result = await api.workflows.runs.listHistory(requestParams);
      } else {
        throw new Error('SDK中没有可用的工作流历史查询方法');
      }
      
      this.logger.log(`获取工作流执行历史成功, 共${result.total || 0}条记录`);
      return result;
    } catch (error) {
      this.logger.error(`获取工作流执行历史失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '获取工作流执行历史失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
  /**
   * 简单上传文件 - 只上传文件并返回file_id
   */
  async uploadFile(file: Buffer, filename: string, mimeType: string = 'application/pdf') {
    try {
      this.logger.log(`上传文件, filename=${filename}, mimeType=${mimeType}, 文件大小=${file.length} bytes`);
      
      // 检查文件是否为空
      if (!file || file.length === 0) {
        throw new Error('文件内容为空');
      }
      
      // 确保client已初始化
      if (!this.client) {
        throw new Error('CozeAPI客户端未初始化');
      }
      
      // 获取token
      try {
        const dynamicClient = this.client as any;
        if (typeof dynamicClient.token === 'function') {
          this.logger.log('获取token...');
          await dynamicClient.token();
          this.logger.log('token获取成功');
        }
      } catch (tokenError) {
        this.logger.warn(`token获取失败: ${tokenError.message}`);
      }
      
      this.logger.log('开始使用直接API方式上传文件...');
      
      // 直接使用fetch API发送请求
      const FormData = require('form-data');
      const formData = new FormData();
      
      // 添加文件到formData
      formData.append('file', file, {
        filename: filename,
        contentType: mimeType,
      });
      
      // 获取token
      let token = this.accessToken;
      if (!token) {
        this.logger.log('重新获取token');
        const tokenResult = await getJWTToken({
          baseURL: COZE_CN_BASE_URL,
          appId: this.defaultConfig.coze_client_id,
          aud: 'api.coze.cn',
          keyid: this.defaultConfig.coze_public_key_id,
          privateKey: this.defaultConfig.coze_private_key,
        });
        token = tokenResult.access_token;
        this.accessToken = token;
      }
      
      // 发送请求
      this.logger.log('发送文件上传请求');

      const fileUploadUrl = `${COZE_CN_BASE_URL}/v1/files/upload`;
      this.logger.log('fileUploadUrl：'+fileUploadUrl);
      this.logger.log('token：'+token);
      
      const response = await axios.post(fileUploadUrl, formData, {
        headers: {
          'Authorization': `Bearer ${token}`,
          ...formData.getHeaders()
        }
      });
      
      this.logger.log(`文件上传响应: ${JSON.stringify(response.data)}`);
      
      // 处理响应
      if (response.data && response.data.code === 0 && response.data.data) {
        // 从嵌套的data对象中获取id
        const fileId = response.data.data.id;
        if (fileId) {
          this.logger.log(`文件上传成功, 返回的fileId=${fileId}`);
          return { file_id: fileId };
        }
      }

      // 如果响应结构不同，尝试其他可能的位置
      if (response.data && (response.data.file_id || response.data.id)) {
        const fileId = response.data.file_id || response.data.id;
        this.logger.log(`文件上传成功, 返回的fileId=${fileId}`);
        return { file_id: fileId };
      }
      
      this.logger.error(`文件上传响应异常: ${JSON.stringify(response.data)}`);
      throw new Error('文件上传响应没有返回文件ID');
    } catch (error) {
      this.logger.error(`上传文件失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      if (error.response) {
        this.logger.error(`错误响应: ${JSON.stringify(error.response.data || {})}`);
      }
      throw new HttpException(
        error.message || '上传文件失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  /**
   * 查询工作流异步执行结果
   * @param workflowId 工作流ID
   * @param executeId 执行ID
   * @param extractField 需要提取的字段，例如"data"
   * @param resultField 结果字段名
   * @returns 工作流执行结果或提取的字段值
   */
  async getWorkflowResult(workflowId: string, executeId: string, extractField?: string, resultField: string = 'result_data') {
    try {
      this.logger.log(`查询工作流执行结果, workflowId=${workflowId}, executeId=${executeId}`);
      
      // 确保client已初始化
      if (!this.client || !(this.client as any).workflows) {
        throw new Error('CozeAPI客户端未初始化或workflows不可用');
      }
      
      // 使用SDK查询工作流执行结果
      const api = this.client as any;
      let result;
      
      // 尝试使用runs.get方法获取结果
      if (api.workflows.runs && typeof api.workflows.runs.get === 'function') {
        this.logger.log('使用runs.get方法获取工作流执行结果');
        result = await api.workflows.runs.get(workflowId, executeId);
      } else if (api.workflows.runs && typeof api.workflows.runs.getResult === 'function') {
        // 尝试其他可能的方法名
        this.logger.log('使用runs.getResult方法获取工作流执行结果');
        result = await api.workflows.runs.getResult(workflowId, executeId);
      } else if (api.workflows.runs && typeof api.workflows.runs.result === 'function') {
        this.logger.log('使用runs.result方法获取工作流执行结果');
        result = await api.workflows.runs.result(workflowId, executeId);
      } else {
        // 如果上述方法都不可用，尝试使用history方法
        this.logger.log('使用history方法获取工作流执行结果');
        result = await this.getWorkflowHistory(workflowId, executeId);
      }
      
      this.logger.log(`获取工作流执行结果成功: ${JSON.stringify(result)}`);
      
      // 用于返回的结果对象
      const responseObj: any = {};
      
      // 检查执行状态
      let isSuccess = false;
      let statusText = '';
      
      // 预设提取字段的默认空值
      if (extractField) {
        responseObj[extractField] = '';
      }
      
      // 始终将原始结果放入指定的结果字段
      responseObj[resultField] = result;
      
      // 获取状态 - 根据返回结果的不同格式进行处理
      if (Array.isArray(result) && result.length > 0) {
        // 从数组结果中获取状态
        const item = result[0];
        
        // 检查状态字段
        if (item.execute_status === 'Success' || item.error_code === '0') {
          isSuccess = true;
          statusText = 'Success';
        } else {
          statusText = item.execute_status || item.error_message || 'Unknown';
        }
        
        // 添加状态到响应
        responseObj.status = statusText;
        responseObj.success = isSuccess;
        
        // 如果执行成功且需要提取特定字段
        if (isSuccess && extractField && item.output) {
          try {
            // 解析output字段
            const outputObj = JSON.parse(item.output);
            
            // 查找Output字段，它可能包含data等信息
            if (outputObj.Output) {
              const outputData = JSON.parse(outputObj.Output);
              
              // 提取特定字段值
              if (extractField === 'data' && outputData.data) {
                responseObj.data = outputData.data;
              } else if (extractField === 'data') {
                responseObj.data = '';  // 明确设置为空字符串
              } else if (outputData[extractField]) {
                responseObj[extractField] = outputData[extractField];
              }
            }
            
            // 也尝试直接从output中查找
            if (outputObj[extractField]) {
              responseObj[extractField] = outputObj[extractField];
            }
          } catch (parseError) {
            this.logger.error(`解析输出字段失败: ${parseError.message}`);
            responseObj.parse_error = parseError.message;
          }
        }
      } else if (result.status) {
        // 处理status字段
        let status = result.status;
        if (typeof status === 'number') {
          // 数字状态码转换为可读状态
          switch(status) {
            case 0: 
              statusText = 'QUEUED'; 
              break;
            case 1: 
              statusText = 'RUNNING'; 
              break;
            case 2: 
              statusText = 'SUCCEEDED'; 
              isSuccess = true;
              break;
            case 3: 
              statusText = 'FAILED'; 
              break;
            case 4: 
              statusText = 'TIMED_OUT'; 
              break;
            default: 
              statusText = `UNKNOWN(${status})`;
          }
        } else {
          statusText = status;
          isSuccess = (status === 'SUCCEEDED');
        }
        
        // 添加状态到响应
        responseObj.status = statusText;
        responseObj.success = isSuccess;
        
        // 如果成功且需要提取字段
        if (isSuccess && extractField && result.output) {
          try {
            // 尝试解析输出
            const outputData = JSON.parse(result.output);
            if (extractField === 'data' && outputData.data) {
              responseObj.data = outputData.data;
            } else if (extractField === 'data') {
              responseObj.data = '';  // 明确设置为空字符串
            } else if (outputData[extractField]) {
              responseObj[extractField] = outputData[extractField];
            }
          } catch (parseError) {
            this.logger.error(`解析输出字段失败: ${parseError.message}`);
            responseObj.parse_error = parseError.message;
          }
        }
      }
      
      // 如果在所有位置都没找到字段值，确保返回空字符串
      if (extractField && !responseObj[extractField]) {
        responseObj[extractField] = '';
      }
      
      return responseObj;
    } catch (error) {
      this.logger.error(`获取工作流执行结果失败: ${error.message}`);
      this.logger.error(`错误堆栈: ${error.stack}`);
      throw new HttpException(
        error.message || '获取工作流执行结果失败',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  // 获取访问令牌的辅助方法
  private async getAccessToken(): Promise<string> {
    if (this.accessToken) {
      return this.accessToken;
    }
    
    try {
      this.logger.log('尝试获取JWT Token');
      const token = await getJWTToken({
        baseURL: COZE_CN_BASE_URL,
        appId: this.defaultConfig.coze_client_id,
        aud: 'api.coze.cn',
        keyid: this.defaultConfig.coze_public_key_id,
        privateKey: this.defaultConfig.coze_private_key,
      });
      
      if (token && token.access_token) {
        this.logger.log(`JWT Token获取成功: ${token.access_token.substring(0, 10)}...`);
        this.accessToken = token.access_token;
        return token.access_token;
      } else {
        this.logger.error('JWT Token获取失败: 返回结果为空或无效');
        throw new Error('Invalid token response');
      }
    } catch (tokenError) {
      this.logger.error(`JWT Token获取失败: ${tokenError.message}`, tokenError.stack);
      throw tokenError;
    }
  }
}