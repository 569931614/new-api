import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  DatePicker,
  Form,
  Input,
  InputNumber,
  Modal,
  Select,
  Switch,
  TextArea,
  Typography,
  Space,
  Divider,
  Tag,
  Banner
} from '@douyinfe/semi-ui';
import { IconShield, IconUser, IconClock, IconCode, IconSetting, IconPlus } from '@douyinfe/semi-icons';
import { API, showError, showSuccess } from '../../helpers';

const { Title } = Typography;

const BatchCreateAuthCode = (props) => {
  const { t } = useTranslation();
  const originInputs = {
    count: 10,
    name: '',
    description: '',
    user_type: 1,
    expired_time: -1,
    is_bot: false,
    wx_auto_x_code: '',
    machine_code: '',
    group: ''
  };

  const [inputs, setInputs] = useState(originInputs);
  const [loading, setLoading] = useState(false);
  const [expiredDate, setExpiredDate] = useState(null);
  const [groupOptions, setGroupOptions] = useState([]);

  const handleInputChange = (name, value) => {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  };

  const fetchGroups = async () => {
    try {
      let res = await API.get(`/api/group/`);
      if (res === undefined) {
        return;
      }
      setGroupOptions(
        res.data.data.map((group) => ({
          label: group,
          value: group,
        })),
      );
    } catch (error) {
      showError(error.message);
    }
  };

  useEffect(() => {
    if (props.visible) {
      fetchGroups();
    }
  }, [props.visible]);

  const submit = async () => {
    if (!inputs.name) {
      showError(t('请填写名称！'));
      return;
    }

    if (inputs.count <= 0 || inputs.count > 100) {
      showError(t('生成数量必须在1-100之间！'));
      return;
    }

    setLoading(true);
    
    // 处理过期时间
    let submitInputs = { ...inputs };
    if (expiredDate) {
      submitInputs.expired_time = Math.floor(expiredDate.getTime() / 1000);
    } else {
      submitInputs.expired_time = -1;
    }

    try {
      let res = await API.post('/api/auth_code/batch', submitInputs);
      const { success, message, data } = res.data;
      if (success) {
        showSuccess(t('批量创建成功！共创建 ') + data.length + t(' 个授权码'));
        setInputs(originInputs);
        setExpiredDate(null);
        props.refresh();
        props.onCancel();
      } else {
        showError(message);
      }
    } catch (error) {
      showError(error.message);
    }
    setLoading(false);
  };

  const userTypeOptions = [
    { label: t('普通用户'), value: 1 },
    { label: t('管理员'), value: 10 },
    { label: t('超级管理员'), value: 100 }
  ];

  return (
    <Modal
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <IconPlus style={{ color: '#1890ff', fontSize: '20px' }} />
          <Title heading={4} style={{ margin: 0 }}>
            {t('批量生成授权码')}
          </Title>
        </div>
      }
      visible={props.visible}
      onCancel={props.onCancel}
      onOk={submit}
      okText={t('生成')}
      cancelText={t('取消')}
      confirmLoading={loading}
      width={600}
      centered={true}
      bodyStyle={{ padding: 0 }}
      style={{ borderRadius: '12px' }}
    >
      <div style={{ padding: '24px' }}>
        <Banner
          type="info"
          description={t('批量生成功能将自动创建指定数量的授权码，每个授权码都是唯一的随机字符串')}
          style={{ marginBottom: '24px' }}
        />

        <Form
          labelPosition={'left'}
          labelAlign={'left'}
          labelWidth={120}
          style={{
            background: 'transparent'
          }}
          initValues={inputs}
        >
          {/* 生成配置区域 */}
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '16px',
              color: '#1890ff',
              fontWeight: 600
            }}>
              <IconCode />
              <span>{t('生成配置')}</span>
            </div>

            <Form.InputNumber
              field={'count'}
              label={
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <span>{t('生成数量')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
              }
              placeholder={t('请输入生成数量')}
              onChange={(value) => handleInputChange('count', value)}
              min={1}
              max={100}
              required
              style={{ marginBottom: '16px' }}
              suffix={<span style={{ color: '#8c8c8c' }}>{t('个')}</span>}
            />
            <div style={{
              fontSize: '12px',
              color: '#8c8c8c',
              marginLeft: '120px',
              marginTop: '-12px',
              marginBottom: '16px'
            }}>
              {t('最多可生成100个授权码')}
            </div>

            <Form.Input
              field={'name'}
              label={
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <span>{t('名称前缀')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
              }
              placeholder={t('请输入名称前缀，将自动添加序号')}
              onChange={(value) => handleInputChange('name', value)}
              required
              style={{ marginBottom: '16px' }}
            />
            <div style={{
              fontSize: '12px',
              color: '#8c8c8c',
              marginLeft: '120px',
              marginTop: '-12px',
              marginBottom: '16px'
            }}>
              {t('例如：测试授权码_1, 测试授权码_2...')}
            </div>

            <Form.TextArea
              field={'description'}
              label={t('描述')}
              placeholder={t('请输入描述信息（可选）')}
              onChange={(value) => handleInputChange('description', value)}
              autosize={{ minRows: 2, maxRows: 4 }}
              style={{ marginBottom: '16px' }}
            />
          </div>

          <Divider margin="16px" />
          {/* 权限配置区域 */}
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '16px',
              color: '#52c41a',
              fontWeight: 600
            }}>
              <IconUser />
              <span>{t('权限配置')}</span>
            </div>

            <Form.Select
              field={'user_type'}
              label={
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <span>{t('用户类型')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
              }
              onChange={(value) => handleInputChange('user_type', value)}
              optionList={userTypeOptions}
              required
              style={{ marginBottom: '16px' }}
            />

            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '16px',
              marginBottom: '16px'
            }}>
              <Form.Switch
                field={'is_bot'}
                label={t('机器人账户')}
                onChange={(checked) => handleInputChange('is_bot', checked)}
              />
              <span style={{
                fontSize: '12px',
                color: '#8c8c8c',
                marginLeft: '8px'
              }}>
                {t('启用后所有授权码将用于机器人账户')}
              </span>
            </div>
          </div>

          <Divider margin="16px" />

          {/* 时间配置区域 */}
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '16px',
              color: '#fa8c16',
              fontWeight: 600
            }}>
              <IconClock />
              <span>{t('时间配置')}</span>
            </div>

            <Form.DatePicker
              field={'expired_time'}
              label={t('过期时间')}
              placeholder={t('选择过期时间，不选择表示永不过期')}
              value={expiredDate}
              onChange={(date) => setExpiredDate(date)}
              type="dateTime"
              style={{ marginBottom: '16px' }}
            />
            <div style={{
              fontSize: '12px',
              color: '#8c8c8c',
              marginLeft: '120px',
              marginTop: '-12px'
            }}>
              {t('不设置过期时间表示永久有效')}
            </div>
          </div>

          <Divider margin="16px" />

          {/* 高级配置区域 */}
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '16px',
              color: '#722ed1',
              fontWeight: 600
            }}>
              <IconSetting />
              <span>{t('高级配置')}</span>
            </div>

            <Form.Input
              field={'wx_auto_x_code'}
              label={t('WxAutoX码')}
              placeholder={t('请输入WxAutoX码（可选）')}
              onChange={(value) => handleInputChange('wx_auto_x_code', value)}
              style={{ marginBottom: '16px' }}
            />
            <div style={{
              fontSize: '12px',
              color: '#8c8c8c',
              marginLeft: '120px',
              marginTop: '-12px',
              marginBottom: '16px'
            }}>
              {t('所有生成的授权码将使用相同的WxAutoX码')}
            </div>

            <Form.Select
              field={'group'}
              label={t('分组')}
              placeholder={t('请选择分组（可选）')}
              value={inputs.group ? inputs.group.split(',').filter(g => g.trim() !== '') : []}
              onChange={(value) => {
                const groupStr = Array.isArray(value) ? value.join(',') : '';
                handleInputChange('group', groupStr);
              }}
              optionList={groupOptions}
              multiple
              allowAdditions
              additionLabel={t('请在系统设置页面编辑分组倍率以添加新的分组：')}
              style={{ marginBottom: '16px' }}
              showClear
            />
            <div style={{
              fontSize: '12px',
              color: '#8c8c8c',
              marginLeft: '120px',
              marginTop: '-12px'
            }}>
              {t('设置分组后，可通过外部接口获取该分组下的渠道列表')}
            </div>
          </div>
        </Form>
      </div>
    </Modal>
  );
};

export default BatchCreateAuthCode;
