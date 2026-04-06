<template>
  <div class="container">
    <el-alert class="info-alert" type="success">
      <p>账户ID: {{ currentAccountId }}</p>
      <p>用户名: {{ currentAccountName }}</p>
      <p>角色: {{ roleText }}</p>
      <p v-if="isDoctor">待审批: {{ pendingList.length }} | 已通过: {{ approvedList.length }} | 已拒绝: {{ rejectedList.length }}</p>
      <p v-if="isPatient">待审批: {{ pendingList.length }} | 已审批: {{ processedList.length }}</p>
    </el-alert>

    <!-- 管理员端：查看已认证的组织 -->
    <div v-if="isAdmin">
      <div class="admin-header">
        <h2>
          <i class="el-icon-medal"></i>
          区块链组织证书管理
        </h2>
        <p class="header-desc">查看联盟链中所有组织及其证书状态（已认证: {{ certifiedOrgs.length }} | 审核中: {{ pendingOrgs.length }}）</p>
      </div>

      <el-card class="summary-card">
        <div class="summary-content">
          <div class="summary-item">
            <div class="summary-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <i class="el-icon-office-building"></i>
            </div>
            <div class="summary-info">
              <p class="summary-value">{{ organizations.length }}</p>
              <p class="summary-label">认证组织</p>
            </div>
          </div>
          <div class="summary-item">
            <div class="summary-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
              <i class="el-icon-connection"></i>
            </div>
            <div class="summary-info">
              <p class="summary-value">{{ totalNodes }}</p>
              <p class="summary-label">区块链节点</p>
            </div>
          </div>
          <div class="summary-item">
            <div class="summary-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
              <i class="el-icon-circle-check"></i>
            </div>
            <div class="summary-info">
              <p class="summary-value">{{ certificationRate }}%</p>
              <p class="summary-label">认证率</p>
            </div>
          </div>
          <div class="summary-item">
            <div class="summary-icon" style="background: linear-gradient(135deg, #ffa726 0%, #fb8c00 100%);">
              <i class="el-icon-loading"></i>
            </div>
            <div class="summary-info">
              <p class="summary-value">{{ pendingOrgs.length }}</p>
              <p class="summary-label">审核中</p>
            </div>
          </div>
        </div>
      </el-card>

      <el-row :gutter="20">
        <el-col v-for="(org, index) in organizations" :key="index" :xs="24" :sm="12" :md="8">
          <el-card class="org-card" shadow="hover">
            <div class="org-header">
              <div class="org-icon" :style="{ background: org.color }">
                <i :class="org.icon"></i>
              </div>
              <div class="org-title">
                <h3>{{ org.name }}</h3>
                <el-tag v-if="org.status === 'pending'" type="warning" size="small">
                  <i class="el-icon-loading"></i>
                  审核中
                </el-tag>
                <el-tag v-else type="success" size="small">
                  <i class="el-icon-circle-check"></i>
                  已认证
                </el-tag>
              </div>
            </div>

            <div class="org-info">
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-postcard"></i>
                  MSP ID
                </span>
                <span class="info-value">{{ org.mspId }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-link"></i>
                  域名
                </span>
                <span class="info-value">{{ org.domain }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-connection"></i>
                  节点数量
                </span>
                <span class="info-value">{{ org.nodeCount }} 个</span>
              </div>
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-position"></i>
                  端口
                </span>
                <span class="info-value">{{ org.ports }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-s-flag"></i>
                  角色
                </span>
                <span class="info-value">{{ org.role }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">
                  <i class="el-icon-time"></i>
                  证书有效期
                </span>
                <span class="info-value">{{ org.certValidity }}</span>
              </div>
            </div>

            <div class="org-footer">
              <el-button type="text" size="small" @click="viewCertDetails(org)">
                <i class="el-icon-view"></i>
                查看证书详情
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 证书详情对话框 -->
      <el-dialog
        title="证书详情"
        :visible.sync="certDialogVisible"
        width="600px"
        :close-on-click-modal="false"
      >
        <div v-if="selectedOrg" class="cert-details">
          <el-alert type="success" :closable="false" style="margin-bottom: 20px;">
            <span>该组织已通过 Hyperledger Fabric MSP 认证，拥有合法的区块链网络访问权限</span>
          </el-alert>

          <div class="cert-section">
            <h4><i class="el-icon-office-building"></i> 组织信息</h4>
            <p><strong>组织名称：</strong>{{ selectedOrg.name }}</p>
            <p><strong>MSP ID：</strong>{{ selectedOrg.mspId }}</p>
            <p><strong>域名：</strong>{{ selectedOrg.domain }}</p>
            <p><strong>角色：</strong>{{ selectedOrg.role }}</p>
          </div>

          <div class="cert-section">
            <h4><i class="el-icon-medal"></i> 证书信息</h4>
            <p><strong>证书类型：</strong>X.509 数字证书</p>
            <p><strong>证书状态：</strong><el-tag type="success" size="small">有效</el-tag></p>
            <p><strong>颁发机构：</strong>Fabric CA</p>
            <p><strong>有效期：</strong>{{ selectedOrg.certValidity }}</p>
          </div>

          <div class="cert-section">
            <h4><i class="el-icon-connection"></i> 节点信息</h4>
            <p><strong>节点数量：</strong>{{ selectedOrg.nodeCount }} 个</p>
            <p><strong>节点端口：</strong>{{ selectedOrg.ports }}</p>
            <p><strong>节点类型：</strong>{{ selectedOrg.nodeType }}</p>
          </div>

          <div class="cert-section">
            <h4><i class="el-icon-key"></i> 权限信息</h4>
            <ul class="permission-list">
              <li v-for="(perm, idx) in selectedOrg.permissions" :key="idx">
                <i class="el-icon-check"></i>
                {{ perm }}
              </li>
            </ul>
          </div>
        </div>
      </el-dialog>
    </div>

    <!-- 病人端：审批授权申请 -->
    <div v-if="isPatient">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <!-- 待审批 -->
        <el-tab-pane label="待审批" name="pending">
          <div v-if="pendingList.length === 0" style="text-align: center;">
            <el-alert title="暂无待审批的授权申请" type="info" />
          </div>

          <el-row v-loading="loading" :gutter="20">
            <el-col v-for="(val, index) in pendingList" :key="index" :xs="24" :sm="12" :md="8">
              <el-card class="request-card" shadow="hover">
                <div class="card-header">
                  <el-tag type="warning" size="small">待审批</el-tag>
                  <span class="request-time">{{ val.request_time }}</span>
                </div>

                <div class="request-info">
                  <div class="info-row">
                    <i class="el-icon-user-solid"></i>
                    <span><strong>申请医生：</strong>{{ val.doctor_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-office-building"></i>
                    <span><strong>所属医院：</strong>{{ val.doctor_org_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-document"></i>
                    <span><strong>申请理由：</strong></span>
                  </div>
                  <div class="reason-text">{{ val.reason }}</div>
                </div>

                <div class="card-actions">
                  <el-button type="success" size="small" @click="handleApprove(val, true)">
                    同意
                  </el-button>
                  <el-button type="danger" size="small" @click="openRejectDialog(val)">
                    拒绝
                  </el-button>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <!-- 已审批 -->
        <el-tab-pane label="已审批" name="processed">
          <div v-if="processedList.length === 0" style="text-align: center;">
            <el-alert title="暂无已审批的记录" type="info" />
          </div>

          <el-row v-loading="loading" :gutter="20">
            <el-col v-for="(val, index) in processedList" :key="index" :xs="24" :sm="12" :md="8">
              <el-card class="request-card" shadow="hover">
                <div class="card-header">
                  <el-tag v-if="val.status === 'approved'" type="success" size="small">已同意</el-tag>
                  <el-tag v-else type="danger" size="small">已拒绝</el-tag>
                  <span class="request-time">{{ val.response_time }}</span>
                </div>

                <div class="request-info">
                  <div class="info-row">
                    <i class="el-icon-user-solid"></i>
                    <span><strong>申请医生：</strong>{{ val.doctor_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-office-building"></i>
                    <span><strong>所属医院：</strong>{{ val.doctor_org_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-document"></i>
                    <span><strong>申请理由：</strong></span>
                  </div>
                  <div class="reason-text">{{ val.reason }}</div>
                  <div v-if="val.status === 'rejected'" class="info-row reject-reason">
                    <i class="el-icon-warning"></i>
                    <span><strong>拒绝理由：</strong>{{ val.reason }}</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 医生端：查看授权申请和已授权病历 -->
    <div v-if="isDoctor">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <!-- 待审批 -->
        <el-tab-pane label="待审批" name="pending">
          <div v-if="pendingList.length === 0" style="text-align: center;">
            <el-alert title="暂无待审批的授权申请" type="info" />
          </div>

          <el-row v-loading="loading" :gutter="20">
            <el-col v-for="(val, index) in pendingList" :key="index" :xs="24" :sm="12" :md="8">
              <el-card class="request-card" shadow="hover">
                <div class="card-header">
                  <el-tag type="warning" size="small">待审批</el-tag>
                  <span class="request-time">{{ val.request_time }}</span>
                </div>

                <div class="request-info">
                  <div class="info-row">
                    <i class="el-icon-user"></i>
                    <span><strong>患者：</strong>{{ val.patient_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-document"></i>
                    <span><strong>申请理由：</strong></span>
                  </div>
                  <div class="reason-text">{{ val.reason }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <!-- 已通过 -->
        <el-tab-pane label="已通过" name="approved">
          <div v-if="approvedList.length === 0" style="text-align: center;">
            <el-alert title="暂无已通过的授权" type="info" />
          </div>

          <el-row v-loading="loading" :gutter="20">
            <el-col v-for="(val, index) in approvedList" :key="index" :xs="24" :sm="12" :md="8">
              <el-card class="request-card approved-card" shadow="hover">
                <div class="card-header">
                  <el-tag type="success" size="small">已授权</el-tag>
                  <span class="request-time">{{ val.response_time }}</span>
                </div>

                <div class="request-info">
                  <div class="info-row">
                    <i class="el-icon-user"></i>
                    <span><strong>患者：</strong>{{ val.patient_name }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-document"></i>
                    <span><strong>病历ID：</strong>{{ val.prescription_id }}</span>
                  </div>
                  <div class="info-row">
                    <i class="el-icon-time"></i>
                    <span><strong>授权时间：</strong>{{ val.response_time }}</span>
                  </div>
                </div>

                <div class="card-actions">
                  <el-button type="primary" size="small" @click="viewPrescription(val.prescription_id)">
                    查看病历
                  </el-button>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <!-- 已拒绝 -->
        <el-tab-pane label="已拒绝" name="rejected">
          <div v-if="rejectedList.length === 0" style="text-align: center;">
            <el-alert title="暂无被拒绝的申请" type="info" />
          </div>

          <el-row v-loading="loading" :gutter="20">
            <el-col v-for="(val, index) in rejectedList" :key="index" :xs="24" :sm="12" :md="8">
              <el-card class="request-card" shadow="hover">
                <div class="card-header">
                  <el-tag type="danger" size="small">已拒绝</el-tag>
                  <span class="request-time">{{ val.response_time }}</span>
                </div>

                <div class="request-info">
                  <div class="info-row">
                    <i class="el-icon-user"></i>
                    <span><strong>患者：</strong>{{ val.patient_name }}</span>
                  </div>
                  <div class="info-row reject-reason">
                    <i class="el-icon-warning"></i>
                    <span><strong>拒绝理由：</strong>{{ val.reason || '未提供' }}</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 拒绝理由对话框（病人端） -->
    <el-dialog
      title="拒绝授权"
      :visible.sync="rejectDialogVisible"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="rejectForm" label-width="100px">
        <el-form-item label="申请医生">
          <el-input v-model="rejectForm.doctorName" disabled />
        </el-form-item>
        <el-form-item label="拒绝理由" required>
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请说明拒绝理由"
          />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="submitReject">确认拒绝</el-button>
      </div>
    </el-dialog>

    <!-- 病历详情对话框（医生端） -->
    <el-dialog
      title="病历详情"
      :visible.sync="prescriptionDialogVisible"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedPrescription" class="detail-content">
        <el-alert type="success" :closable="false" style="margin-bottom: 20px;">
          <span>此病历已获得患者授权，您可以查看完整内容</span>
        </el-alert>

        <!-- 基本信息 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-document"></i>
            <span>基本信息</span>
          </div>
          <div class="section-text">
            <p><strong>病历编号：</strong>{{ selectedPrescription.prescription_no }}</p>
            <p><strong>就诊时间：</strong>{{ selectedPrescription.created }}</p>
            <p><strong>就诊医院：</strong>{{ selectedPrescription.hospital_name }}</p>
            <p><strong>科室：</strong>{{ selectedPrescription.department }}</p>
          </div>
        </div>

        <!-- 患者信息 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-user"></i>
            <span>患者信息</span>
          </div>
          <div class="section-text">
            <p><strong>姓名：</strong>{{ selectedPrescription.patient_name }}</p>
          </div>
        </div>

        <!-- 医生信息 -->
        <div class="info-section">
          <div class="section-header">
            <i class="el-icon-user-solid"></i>
            <span>医生信息</span>
          </div>
          <div class="section-text">
            <p><strong>医生：</strong>{{ selectedPrescription.doctor_name }}</p>
            <p><strong>职称：</strong>{{ selectedPrescription.doctor_title }}</p>
          </div>
        </div>

        <!-- 主诉 -->
        <div v-if="selectedPrescription.chief_complaint" class="info-section">
          <div class="section-header">
            <i class="el-icon-chat-line-square"></i>
            <span>主诉</span>
          </div>
          <div class="section-text">
            <p>{{ selectedPrescription.chief_complaint }}</p>
          </div>
        </div>

        <!-- 现病史 -->
        <div v-if="selectedPrescription.present_illness" class="info-section">
          <div class="section-header">
            <i class="el-icon-document-copy"></i>
            <span>现病史</span>
          </div>
          <div class="section-text">
            <p>{{ selectedPrescription.present_illness }}</p>
          </div>
        </div>

        <!-- 体格检查 -->
        <div v-if="selectedPrescription.physical_exam" class="info-section">
          <div class="section-header">
            <i class="el-icon-view"></i>
            <span>体格检查</span>
          </div>
          <div class="section-text">
            <p>{{ selectedPrescription.physical_exam }}</p>
          </div>
        </div>

        <!-- 诊断结果 -->
        <div class="diagnosis-highlight">
          <div class="diagnosis-header">
            <i class="el-icon-document-checked"></i>
            <span>诊断结果</span>
          </div>
          <div class="diagnosis-text">{{ selectedPrescription.diagnosis }}</div>
        </div>

        <!-- 处方药品 -->
        <div v-if="selectedPrescription.drug && selectedPrescription.drug.length > 0" class="drug-highlight">
          <div class="drug-header">
            <i class="el-icon-medicine-box"></i>
            <span>处方药品</span>
          </div>
          <div class="drug-list-compact">
            <div v-for="(drug, idx) in selectedPrescription.drug" :key="idx" class="drug-item-compact">
              <span class="drug-name">{{ drug.Name }}</span>
              <span class="drug-amount">× {{ drug.amount }}</span>
            </div>
          </div>
        </div>

        <!-- 医嘱 -->
        <div v-if="selectedPrescription.medical_advice" class="info-section">
          <div class="section-header">
            <i class="el-icon-document-add"></i>
            <span>医嘱</span>
          </div>
          <div class="section-text">
            <p>{{ selectedPrescription.medical_advice }}</p>
          </div>
        </div>

        <!-- 备注 -->
        <div v-if="selectedPrescription.comment" class="info-section">
          <div class="section-header">
            <i class="el-icon-edit"></i>
            <span>备注</span>
          </div>
          <div class="section-text">
            <p>{{ selectedPrescription.comment }}</p>
          </div>
        </div>
      </div>
      <div v-else style="text-align: center; padding: 40px;">
        <i class="el-icon-loading" style="font-size: 32px;"></i>
        <p>加载中...</p>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { queryAccessRequests, approveAccess } from '@/api/access'
import { queryPrescriptionList } from '@/api/prescription'
import { mapGetters } from 'vuex'

export default {
  name: 'Authorization',
  data() {
    return {
      activeTab: 'pending',
      loading: false,
      pendingList: [],
      processedList: [],
      approvedList: [],
      rejectedList: [],
      rejectDialogVisible: false,
      rejectForm: {
        requestId: '',
        doctorName: '',
        reason: ''
      },
      prescriptionDialogVisible: false,
      selectedPrescription: null,
      currentAccountId: '',
      currentAccountName: '',
      // 管理员端：组织证书信息
      organizations: [
        {
          name: '协和医院',
          mspId: 'TaobaoMSP',
          domain: 'taobao.com',
          nodeCount: 2,
          ports: '7051, 17051',
          role: '医疗机构',
          certValidity: '2026-02-20 至 2036-02-18（10年）',
          nodeType: 'Peer节点（背书节点 + 提交节点）',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          permissions: [
            '参与交易背书',
            '提交交易到账本',
            '访问通道数据',
            '执行链码',
            '查询账本状态'
          ]
        },
        {
          name: '301医院',
          mspId: 'JDMSP',
          domain: 'jd.com',
          nodeCount: 2,
          ports: '27051, 37051',
          role: '医疗机构',
          certValidity: '2026-02-20 至 2036-02-18（10年）',
          nodeType: 'Peer节点（背书节点 + 提交节点）',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          permissions: [
            '参与交易背书',
            '提交交易到账本',
            '访问通道数据',
            '执行链码',
            '查询账本状态'
          ]
        },
        {
          name: '温江医疗中心',
          mspId: 'WenjinMSP',
          domain: 'wenjin.com',
          nodeCount: 2,
          ports: '47051, 57051',
          role: '医疗机构',
          certValidity: '2026-02-20 至 2036-02-18（10年）',
          nodeType: 'Peer节点（背书节点 + 提交节点）',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
          permissions: [
            '参与交易背书',
            '提交交易到账本',
            '访问通道数据',
            '执行链码',
            '查询账本状态'
          ]
        },
        {
          name: '监管中心',
          mspId: 'RegCenterMSP',
          domain: 'regcenter.com',
          nodeCount: 2,
          ports: '8051, 9051',
          role: '监管机构',
          certValidity: '2026-02-20 至 2036-02-18（10年）',
          nodeType: 'Peer节点（背书节点 + 提交节点）',
          icon: 'el-icon-s-custom',
          color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          permissions: [
            '参与交易背书',
            '提交交易到账本',
            '访问通道数据',
            '执行链码',
            '查询账本状态',
            '监管所有数据'
          ]
        },
        {
          name: '排序服务',
          mspId: 'OrdererMSP',
          domain: 'qq.com',
          nodeCount: 1,
          ports: '7050',
          role: '排序节点',
          certValidity: '2026-02-20 至 2036-02-18（10年）',
          nodeType: 'Orderer节点',
          icon: 'el-icon-sort',
          color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          permissions: [
            '对交易进行排序',
            '生成区块',
            '广播区块到通道',
            '维护通道配置'
          ]
        },
        {
          name: '温江人寿',
          mspId: 'InsuranceMSP',
          domain: 'insurance.com',
          nodeCount: 2,
          ports: '待分配',
          role: '保险机构',
          certValidity: '审核中',
          nodeType: 'Peer节点（待审核）',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #ffa726 0%, #fb8c00 100%)',
          status: 'pending',
          permissions: [
            '待审核通过后分配权限'
          ]
        },
        {
          name: '成都中医药大学第二附属医院',
          mspId: 'CDUTCMMSP',
          domain: 'cdutcm.edu.cn',
          nodeCount: 2,
          ports: '待分配',
          role: '医疗机构',
          certValidity: '审核中',
          nodeType: 'Peer节点（待审核）',
          icon: 'el-icon-office-building',
          color: 'linear-gradient(135deg, #ff7043 0%, #f4511e 100%)',
          status: 'pending',
          permissions: [
            '待审核通过后分配权限'
          ]
        },
        {
          name: '四川省药品监督管理局',
          mspId: 'DrugRegMSP',
          domain: 'scda.gov.cn',
          nodeCount: 2,
          ports: '待分配',
          role: '药品监管机构',
          certValidity: '审核中',
          nodeType: 'Peer节点（待审核）',
          icon: 'el-icon-s-check',
          color: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
          status: 'pending',
          permissions: [
            '待审核通过后分配权限'
          ]
        }
      ],
      certDialogVisible: false,
      selectedOrg: null
    }
  },
  computed: {
    ...mapGetters(['roles']),
    isPatient() {
      return this.roles && this.roles[0] === 'patient'
    },
    isDoctor() {
      return this.roles && this.roles[0] === 'doctor'
    },
    isAdmin() {
      return this.roles && this.roles[0] === 'admin'
    },
    roleText() {
      const roleMap = {
        'patient': '病人',
        'doctor': '医生',
        'admin': '管理员'
      }
      return roleMap[this.roles[0]] || '用户'
    },
    totalNodes() {
      return this.organizations.reduce((sum, org) => sum + org.nodeCount, 0)
    },
    certifiedOrgs() {
      return this.organizations.filter(org => org.status !== 'pending')
    },
    pendingOrgs() {
      return this.organizations.filter(org => org.status === 'pending')
    },
    certificationRate() {
      const certified = this.certifiedOrgs.length
      const total = this.organizations.length
      return total > 0 ? Math.round((certified / total) * 100) : 0
    }
  },
  mounted() {
    // 从sessionStorage获取用户信息
    const userInfo = JSON.parse(sessionStorage.getItem('userInfo'))
    if (userInfo && userInfo.account_id) {
      this.currentAccountId = userInfo.account_id
      this.currentAccountName = userInfo.account_name
      this.loadRequests()
    } else {
      this.$message.error('无法获取用户信息，请重新登录')
    }
  },
  methods: {
    handleTabClick() {
      this.loadRequests()
    },
    // 加载授权请求列表（带静默重试机制）
    async loadRequests(retryCount = 0) {
      const maxRetries = 3
      
      if (!this.currentAccountId) {
        this.$message.error('无法获取用户ID，请重新登录')
        return
      }
      
      // 只在第一次加载时显示loading
      if (retryCount === 0) {
        this.loading = true
      }
      
      const role = this.isPatient ? 'patient' : 'doctor'
      
      try {
        const res = await queryAccessRequests({
          user_id: this.currentAccountId,
          role: role
        })
        
        this.loading = false
        
        if (res.code === 200) {
          const allRequests = res.data || []
          
          if (this.isPatient) {
            // 病人端：分为待审批和已审批
            this.pendingList = allRequests.filter(r => r.status === 'pending')
            this.processedList = allRequests.filter(r => r.status !== 'pending')
          } else {
            // 医生端：分为待审批、已通过、已拒绝
            this.pendingList = allRequests.filter(r => r.status === 'pending')
            this.approvedList = allRequests.filter(r => r.status === 'approved')
            this.rejectedList = allRequests.filter(r => r.status === 'rejected')
          }
        } else {
          // 如果失败且还有重试次数，静默重试
          if (retryCount < maxRetries) {
            console.log(`查询授权请求失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return this.loadRequests(retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            this.$message.error(res.msg || '查询失败')
          }
        }
      } catch (err) {
        // 只在控制台记录错误
        if (retryCount === 0) {
          console.error('查询授权请求失败:', err)
        }
        
        // 如果失败且还有重试次数，静默重试
        if (retryCount < maxRetries) {
          console.log(`查询授权请求失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.loadRequests(retryCount + 1)
        } else {
          // 所有重试都失败后才显示错误
          this.loading = false
          this.$message.error('查询失败：' + err.message)
        }
      }
    },
    // 处理审批（带静默重试机制）
    async handleApprove(request, approved, retryCount = 0) {
      const maxRetries = 3
      const action = approved ? '同意' : '拒绝'
      
      try {
        await this.$confirm(`确认${action}该授权申请吗？`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        try {
          const res = await approveAccess({
            request_id: request.id,
            patient_id: this.currentAccountId,
            approved: approved ? 'true' : 'false',
            reject_reason: ''
          })
          
          if (res.code === 200) {
            this.$message.success(`已${action}授权申请`)
            this.loadRequests()
          } else {
            // 如果失败且还有重试次数，静默重试
            if (retryCount < maxRetries) {
              console.log(`${action}授权失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
              await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
              return this.handleApprove(request, approved, retryCount + 1)
            } else {
              // 所有重试都失败后才显示错误
              this.$message.error(res.msg || `${action}失败`)
            }
          }
        } catch (err) {
          // 只在控制台记录错误
          if (retryCount === 0) {
            console.error(`${action}授权失败:`, err)
          }
          
          // 如果失败且还有重试次数，静默重试
          if (retryCount < maxRetries) {
            console.log(`${action}授权失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return this.handleApprove(request, approved, retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            this.$message.error(`${action}失败：` + err.message)
          }
        }
      } catch {
        // 用户取消操作
      }
    },
    openRejectDialog(request) {
      this.rejectForm = {
        requestId: request.id,
        doctorName: request.doctor_name,
        reason: ''
      }
      this.rejectDialogVisible = true
    },
    // 提交拒绝（带静默重试机制）
    async submitReject(retryCount = 0) {
      const maxRetries = 3
      
      if (!this.rejectForm.reason.trim()) {
        this.$message.warning('请填写拒绝理由')
        return
      }

      try {
        const res = await approveAccess({
          request_id: this.rejectForm.requestId,
          patient_id: this.currentAccountId,
          approved: 'false',
          reject_reason: this.rejectForm.reason
        })
        
        if (res.code === 200) {
          this.$message.success('已拒绝授权申请')
          this.rejectDialogVisible = false
          this.loadRequests()
        } else {
          // 如果失败且还有重试次数，静默重试
          if (retryCount < maxRetries) {
            console.log(`拒绝授权失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
            return this.submitReject(retryCount + 1)
          } else {
            // 所有重试都失败后才显示错误
            this.$message.error(res.msg || '拒绝失败')
          }
        }
      } catch (err) {
        // 只在控制台记录错误
        if (retryCount === 0) {
          console.error('拒绝授权失败:', err)
        }
        
        // 如果失败且还有重试次数，静默重试
        if (retryCount < maxRetries) {
          console.log(`拒绝授权失败，正在静默重试 (${retryCount + 1}/${maxRetries})...`)
          await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1)))
          return this.submitReject(retryCount + 1)
        } else {
          // 所有重试都失败后才显示错误
          this.$message.error('拒绝失败：' + err.message)
        }
      }
    },
    // 医生端：查看已授权的病历
    viewPrescription(prescriptionId) {
      // 跳转到病历详情页面，可以添加补充记录
      this.$router.push({
        name: 'PrescriptionDetail',
        query: { id: prescriptionId }
      })
    },
    // 管理员端：查看证书详情
    viewCertDetails(org) {
      this.selectedOrg = org
      this.certDialogVisible = true
    }
  }
}
</script>

<style scoped>
.container {
  padding: 20px;
}

/* 管理员端样式 */
.admin-header {
  margin-bottom: 30px;
  text-align: center;
}

.admin-header h2 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.admin-header h2 i {
  color: #4facfe;
}

.header-desc {
  color: #909399;
  font-size: 14px;
}

.summary-card {
  margin-bottom: 30px;
  border-radius: 12px;
}

.summary-content {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  gap: 20px;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.summary-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.summary-icon i {
  font-size: 28px;
  color: white;
}

.summary-info {
  text-align: left;
}

.summary-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  margin: 0;
  line-height: 1;
}

.summary-label {
  font-size: 14px;
  color: #909399;
  margin: 5px 0 0 0;
}

.org-card {
  margin-bottom: 20px;
  border-radius: 12px;
  transition: all 0.3s;
}

.org-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.org-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 2px solid #f5f7fa;
}

.org-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.org-icon i {
  font-size: 24px;
  color: white;
}

.org-title {
  flex: 1;
}

.org-title h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #303133;
}

.org-info {
  margin-bottom: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f5f7fa;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

.info-label i {
  color: #409EFF;
}

.info-value {
  color: #303133;
  font-weight: 500;
  font-size: 14px;
}

.org-footer {
  padding-top: 15px;
  border-top: 1px solid #f5f7fa;
  text-align: center;
}

/* 证书详情对话框 */
.cert-details {
  max-height: 500px;
  overflow-y: auto;
}

.cert-section {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.cert-section h4 {
  margin: 0 0 15px 0;
  font-size: 16px;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.cert-section h4 i {
  color: #409EFF;
}

.cert-section p {
  margin: 8px 0;
  color: #606266;
  line-height: 1.8;
}

.permission-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.permission-list li {
  padding: 8px 0;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 8px;
}

.permission-list li i {
  color: #67C23A;
  font-weight: bold;
}

/* 原有样式 */
.info-alert {
  margin-bottom: 20px;
}

.request-card {
  margin-bottom: 20px;
  transition: all 0.3s;
}

.request-card:hover {
  transform: translateY(-3px);
}

.approved-card {
  border-left: 4px solid #67C23A;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #EBEEF5;
}

.request-time {
  font-size: 12px;
  color: #909399;
}

.request-info {
  margin-bottom: 15px;
}

.info-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
  color: #606266;
}

.info-row i {
  margin-right: 8px;
  margin-top: 3px;
  color: #409EFF;
}

.reason-text {
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 10px;
}

.reject-reason {
  color: #F56C6C;
  margin-top: 10px;
}

.card-actions {
  display: flex;
  justify-content: center;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid #EBEEF5;
}

/* 病历详情样式 */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.info-section {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.section-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.section-header i {
  margin-right: 8px;
  color: #409EFF;
}

.section-text {
  color: #606266;
  line-height: 1.8;
}

.section-text p {
  margin: 5px 0;
}

.diagnosis-highlight {
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.diagnosis-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 12px;
}

.diagnosis-header i {
  margin-right: 8px;
}

.diagnosis-text {
  font-size: 16px;
  line-height: 1.8;
}

.drug-highlight {
  margin-bottom: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  border-radius: 8px;
  color: white;
}

.drug-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 12px;
}

.drug-header i {
  margin-right: 8px;
}

.drug-list-compact {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.drug-item-compact {
  background-color: rgba(255, 255, 255, 0.2);
  padding: 8px 15px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.drug-name {
  font-weight: 500;
}

.drug-amount {
  font-size: 14px;
  opacity: 0.9;
}
</style>
